package controller

import (
	"net/http"
	"encoding/json"
	"beluga/src/configuration/system/controller"
	"beluga/src/configuration/system/drive"
)

type ConfigController struct {
	controller.Controller
	config []drive.Configuration
}

func (c *ConfigController) GetConf() {
	var appid, namespace string

	req_get_data := c.Request.URL.Query()
	if req_get_data["appid"] == nil {
		c.Json(http.StatusInternalServerError, "err", "缺少参数", "")
	}
	if req_get_data["namespace"] == nil {
		c.Json(http.StatusInternalServerError, "err", "缺少参数", "")
	}

	appid = req_get_data.Get("appid")
	namespace = req_get_data.Get("namespace")

	var conf_cache interface{}
	redis := c.Redis
	name := appid + "_" + namespace + "_json"

	conf_cache, err := redis.HGet("beluga_configuration", name).Result()

	var data interface{}

	if (conf_cache == "") || (err != nil) {
		data = c.getConfigList(appid, namespace)

		// 值存取redis
		data_json, _ := json.Marshal(data)
		redis.HSet("beluga_configuration", name, data_json)
	} else {
		var undata map[string]interface{}

		json.Unmarshal([]byte(conf_cache.(string)), &undata)
		data = undata
	}

	// 正常业务逻辑处理
	c.Json(http.StatusOK, "OK", "请求成功", data)
}

// 同步配置
func (c *ConfigController) SyncConfig() {
	req_get_data := c.Request.URL.Query()
	if req_get_data["appid"] == nil {
		c.Json(http.StatusInternalServerError, "err", "缺少参数", "")
	}
	if req_get_data["namespace"] == nil {
		c.Json(http.StatusInternalServerError, "err", "缺少参数", "")
	}

	redis := c.Redis
	name := req_get_data.Get("appid") + "_" + req_get_data.Get("namespace") + "_json"
	data := c.getConfigList(req_get_data.Get("appid"), req_get_data.Get("namespace"))

	data_json, _ := json.Marshal(data)
	redis.HSet("beluga_configuration", name, data_json)

	// 下发通知

	c.Json(http.StatusOK, "OK", "成功", "")
}

// 获取配置
func (c *ConfigController) getConfigList(appid, namespace string) map[string]interface{} {
	db := c.Mysql.New()
	var conf []drive.Configuration

	db.Table("beluga_configuration").Select("`key`,`val`").Where("`appid`=? and `namespace_name`=?", appid, namespace).Find(&conf)
	c.config = conf

	return drive.SetKeyToJson(conf)
}

// 配置转为struct格式
func (c *ConfigController) handleConfiguration(conf string) []drive.Configuration {
	var config_data []drive.Configuration
	json.Unmarshal([]byte(conf), &config_data)

	return config_data
}

