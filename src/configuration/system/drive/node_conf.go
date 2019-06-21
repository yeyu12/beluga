package drive

import (
	"beluga/src/beluga/configuration_constant"
	"beluga/src/beluga/library"
	"beluga/src/configuration/system/langurage"
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"encoding/json"
	"strconv"
	"fmt"
)

//{
//"appid":"11f04d90bc2cf8fca0d028deb0cf610b",
//"namespace": {
//"明明空间天啊及的三": "http"
//"明明空间天啊及的三": "php|路径+文件名|Sd"
//"明明空间天啊及的三": "json|路径+文件名"
//}
//}

var conf_str = ""

// 三种类型，1网络请求格式，2json文件格式，3php文件格式，后两种是要写入文件

// 初始化配置
func InitNodeConf() {
	key := configuration_constant.NODE_CONF_DIR + library.G_conf_etcd_client.LocalIp
	resp, err := library.G_conf_etcd_client.Kv.Get(context.TODO(), key)

	if err != nil {
		Err(logrus.Fields{}, errors.Wrap(err, "配置初始化失败！"))
	}

	for _, v := range resp.Kvs {
		var node_conf []configuration_constant.NodeConf
		json.Unmarshal(v.Value, &node_conf)

		for _, val := range node_conf {
			for _, v := range val.Namespace {
				G_node_conf[val.Appid+"_"+v.Name] = v.Path
				setFileWrite(val.Appid, v, SetKeyToJson(getConfMysql(val.Appid, v.Name)))
			}
		}
	}
}

// 监听节点配置变化
func watchNode(event *clientv3.Event) {
	var node_conf []configuration_constant.NodeConf
	json.Unmarshal(event.Kv.Value, &node_conf)

	for _, val := range node_conf {
		for _, v := range val.Namespace {
			G_node_conf[val.Appid+"_"+v.Name] = v.Path
			setFileWrite(val.Appid, v, SetKeyToJson(getConfMysql(val.Appid, v.Name)))
		}
	}
}

// 监听节点配置变化
func watchNodeConf(event *clientv3.Event) {
	node_conf_temp := strings.Split(string(event.Kv.Key), "/")
	appid_namespace := strings.Split(node_conf_temp[4], "_")

	if node_conf_temp[3] != library.G_conf_etcd_client.LocalIp {
		return
	}

	appid := appid_namespace[0]
	namespace := appid_namespace[1]
	fileinfo := configuration_constant.NodeConfNamespace{
		Name: namespace,
		Path: string(event.Kv.Value),
	}

	G_node_conf[appid+"_"+namespace] = fileinfo.Path
	setFileWrite(appid, fileinfo, SetKeyToJson(getConfMysql(appid, namespace)))

	// 删除etcd中的数据
	library.G_conf_etcd_client.Kv.Delete(context.TODO(), string(event.Kv.Key))
}

// 中心配置变化同步
func watchConfigurationSync(event *clientv3.Event) {
	node_conf_temp := strings.Split(string(event.Kv.Key), "/")
	appid_namespace := strings.Split(node_conf_temp[3], "_")

	if G_node_conf[node_conf_temp[3]] == "" {
		return
	}

	appid := appid_namespace[0]
	namespace := strings.Replace(strings.Trim(fmt.Sprint(appid_namespace[1:]), "[]"), " ", "_", -1)
	fileinfo := configuration_constant.NodeConfNamespace{
		Name: namespace,
		Path: G_node_conf[node_conf_temp[3]],
	}
	G_node_conf[appid+"_"+namespace] = fileinfo.Path

	setFileWrite(appid, fileinfo, SetKeyToJson(getConfMysql(appid, namespace)))
}

// 获取配置
func getConfMysql(appid, namespace string) []Configuration {
	db := G_mysql.New()
	var conf []Configuration

	db.Table("beluga_configuration").Select("`key`,`val`").Where("`appid`=? and `namespace_name`=?", appid, namespace).Find(&conf)

	return conf
}

// setFileWrite 文件写入
// appid string
// fileinfo string 文件信息
// data interface 数据
func setFileWrite(appid string, fileinfo configuration_constant.NodeConfNamespace, data interface{}) {
	data_str, _ := json.Marshal(data)
	file_to_type := strings.Split(fileinfo.Path, "|")
	node_conf_type := strings.ToLower(file_to_type[0])

	switch node_conf_type {
	case NODE_CONF_JSON_TYPE:
		if ioutil.WriteFile(file_to_type[1], data_str, 0644) != nil {
			Notices(logrus.Fields{}, "json文件写入失败")
		}
		break
	case NODE_CONF_PHP_TYPE:
		mapToStr(data, "")
		// 判断框架类型，保存为相应的结构
		conf_php_str := langurage.StrToPhpFileFormat(file_to_type[2], conf_str, fileinfo.Name)

		if ioutil.WriteFile(file_to_type[1], []byte(conf_php_str), 0644) != nil {
			Notices(logrus.Fields{}, "php文件写入失败")
		}

		conf_str = ""

		break
	}

	// 直接写入到redis中，对外接口可直接放问
	name := appid + "_" + fileinfo.Name + "_json"
	G_redis.HSet(configuration_constant.CONFIGURATION_REDIS_KEY, name, data_str)
}

// json字符串转php字符串
func mapToStr(data interface{}, str1 string) {
	temp := str1

	if data != nil {
		for key, val := range data.(map[string]interface{}) {
			switch val.(type) {
			case map[string]interface{}:
				str1 += "[\"" + key + "\"]"
				mapToStr(val, str1)
				str1 = temp

				break
			case string:
				val_str := strings.Trim(val.(string), "")

				conf_php_str := ""
				if val_str == "" {
					conf_php_str += str1 + "[\"" + key + "\"]" + "=false;\n"
				} else {
					// TODO 只接受数组对象，和普通字符串
					if val_str[0] == '[' && val_str[len(val_str)-1] == ']' {
						php_str := jsonArrayToPhpStr(val_str)
						conf_php_str += str1 + "[\"" + key + "\"]" + "=" + php_str + ";\n"
					} else {
						conf_php_str += str1 + "[\"" + key + "\"]" + "=\"" + escapeCharacter(val.(string)) + "\";\n"
					}
				}

				conf_str += conf_php_str

				break
			}
		}
	}
}

// json字符串转php字符串
func jsonArrayToPhpStr(conf interface{}) string {
	var conf_json []interface{}

	switch conf.(type) {
	case []interface{}:
		conf_json = conf.([]interface{})
		break
	case string:
		json.Unmarshal([]byte(conf.(string)), &conf_json)
		break
	}

	conf_php_str := "["

	for _, val := range conf_json {
		switch val.(type) {
		case map[string]interface{}:
			conf_php_str += jsonObjToPhpStr(val.(map[string]interface{}))
			break
		case string:
			conf_php_str += "\"" + escapeCharacter(val.(string)) + "\","
			break
		}
	}

	conf_php_str = conf_php_str[:len(conf_php_str)-1]
	conf_php_str += "]"

	return conf_php_str
}

// json字符串转php字符串
func jsonObjToPhpStr(conf map[string]interface{}) string {
	conf_php_str := "["

	for key, val := range conf {
		conf_php_str += "\"" + key + "\""

		switch val.(type) {
		case map[string]interface{}:
			conf_php_str += "=>" + jsonObjToPhpStr(conf) + ","
			break
		case []interface{}:
			conf_php_str += "=>" + jsonArrayToPhpStr(val) + ","
			break
		case string:
			conf_php_str += "=>\"" + escapeCharacter(val.(string)) + "\","
			break
		case int64:
			val_int64_str := strconv.FormatInt(val.(int64), 10)
			conf_php_str += "=>" + val_int64_str + ","
			break
		case int:
			val_int_str := strconv.Itoa(val.(int))
			conf_php_str += "=>" + val_int_str + ","
			break
		case float32:
			val_float32 := val.(float32)
			val_float32_str := strconv.FormatFloat(float64(val_float32), 'f', -1, 32)
			conf_php_str += "=>" + val_float32_str + ","
			break
		case float64:
			val_float64_str := strconv.FormatFloat(val.(float64), 'f', -1, 64)
			conf_php_str += "=>" + val_float64_str + ","
			break
		}
	}

	conf_php_str = conf_php_str[:len(conf_php_str)-1]
	conf_php_str += "],"

	return conf_php_str
}

// 字符转义
func escapeCharacter(str string) string {
	escape_char := []string{
		"$",
	}

	for _, v := range escape_char {
		str = strings.Replace(str, v, "\\"+v, -1)
	}

	return str
}
