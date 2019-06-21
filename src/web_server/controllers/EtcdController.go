package controllers

import (
	"github.com/astaxie/beego"
	"beluga/src/web_server/helpers"
	"encoding/json"
	"strconv"
	"github.com/coreos/etcd/clientv3"
	"time"
	"io/ioutil"
	"beluga/src/beluga/library"
	"strings"
	"fmt"
)

type EtcdController struct {
	BaseController
}

// 获取etcd服务列表
func (c *EtcdController) GetEtcdNode() {
	etcd_ips := beego.AppConfig.DefaultStrings("etcd_host", []string{})

	var etcd_ip_map []map[string]interface{}

	etcd_timeout := time.Duration(beego.AppConfig.DefaultInt("etcd_timeoute", 5000)) * time.Millisecond

	for _, v := range etcd_ips {
		etcd_ip_data := make(map[string]interface{})
		etcd_ip_data["ip"] = v
		etcd_ip_data["status"] = true

		// 判断etcd是否可用
		if !c.verEtcdState(v, etcd_timeout) {
			etcd_ip_data["status"] = false
		}

		etcd_ip_map = append(etcd_ip_map, etcd_ip_data)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, etcd_ip_map)
}

// etcd服务ip地址写入配置
func (c *EtcdController) AddEtcdNodeConf() {
	// 添加etcd服务地址.先检查地址是否可用,如果可用则添加并且重新初始化.最后写入到配置文件中去,下次启动依然有效
	request_data := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	ip := request_data["ip"].(string)
	port := request_data["port"].(float64)

	if (ip == "") || (port == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	ip_port := ip + ":" + strconv.Itoa(int(port))

	// 判断链接是否存在
	etcd_conf_ip := beego.AppConfig.DefaultStrings("etcd_host", []string{})

	for _, v := range etcd_conf_ip {
		if v == ip_port {
			c.ResponseJson(helpers.ETCD_SERVER_EXISTENCE_CODE, false, helpers.ETCD_SERVER_EXISTENCE_MSG)
			return
		}
	}

	etcd_timeout := time.Duration(beego.AppConfig.DefaultInt("etcd_timeoute", 5000)) * time.Millisecond
	res_map := map[string]interface{}{
		"ip":     ip_port,
		"status": true,
	}

	if !c.verEtcdState(ip_port, etcd_timeout) {
		beego.Error("服务不存在")
		c.ResponseJson(helpers.ETCD_SERVER_NOTE_CODE, false, helpers.ETCD_SERVER_NOTE_MSG)
	} else {
		beego.Error("服务存在")
		return
	}

	// 配置写入文件中去
	etcd_ip := beego.AppConfig.DefaultString("etcd_host", "")
	etcd_ip += ";" + ip_port

	c.saveEtcdConf(etcd_ip)
	c.reloadEtcd()

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, res_map)
}

// 删除etcd服务
func (c *EtcdController) DelEtcdNode() {
	request_data := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	ip_port := request_data["ip"].(string)

	if ip_port == "" {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	etcd_conf_ip := beego.AppConfig.DefaultStrings("etcd_host", []string{})

	for k, v := range etcd_conf_ip {
		if v == ip_port {
			etcd_conf_ip = append(etcd_conf_ip[:k], etcd_conf_ip[k+1:]...)
			goto ETCD_SERVER_NOTEEXISTENCE
		}
	}

	c.ResponseJson(helpers.ETCD_SERVER_NOTE_EXISTENCE_CODE, false, helpers.ETCD_SERVER_NOTE_EXISTENCE_MSG)

ETCD_SERVER_NOTEEXISTENCE:
	//删除配置中的

	ip_port_str := strings.Replace(strings.Trim(fmt.Sprint(etcd_conf_ip), "[]"), " ", ";", -1)
	c.saveEtcdConf(ip_port_str)
	c.reloadEtcd()

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 判断etcd状态
func (c *EtcdController) verEtcdState(ip_port string, timeout time.Duration) bool {
	etcd_config := clientv3.Config{
		Endpoints:   []string{ip_port},
		DialTimeout: timeout,
	}

	client, err := clientv3.New(etcd_config)
	if err != nil {
		return false
	}

	defer client.Close()

	return true
}

// 配置文件写入
func (c *EtcdController) saveEtcdConf(etcd_ip string) {
	etcd_conf := "etcd_host=" + etcd_ip + "\n" + "etcd_timeout=" + beego.AppConfig.DefaultString("etcd_timeout", "5000")
	beego.AppConfig.Set("etcd_host", etcd_ip)

	if err := ioutil.WriteFile("./conf/etcd.conf", []byte(etcd_conf), 0644); err != nil {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}
}

// etcd服务重启
func (c *EtcdController)reloadEtcd() {
	library.G_conf_etcd_client.Client.Close()
	library.G_conf_etcd_client.Lease.Close()
	err := library.InitRegister(beego.AppConfig.DefaultStrings("etcd_host", []string{}), beego.AppConfig.DefaultInt("etcd_timeoute", 5000))
	if err != nil {
		c.ResponseJson(helpers.ETCD_SERVER_INIT_CODE, false, helpers.ETCD_SERVER_INIT_MSG)
	}
}
