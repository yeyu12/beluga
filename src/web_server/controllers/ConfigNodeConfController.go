package controllers

import (
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"encoding/json"
	"strconv"
	"time"
	"beluga/src/beluga/configuration_constant"
	"beluga/src/beluga/library"
	"context"
)

type ConfigNodeConfController struct {
	BaseController
}

// 节点配置列表
func (c *ConfigNodeConfController) NodeConfList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]

	if page == 0 {
		page = 1
	}

	configuration_node_conf_model := models.NewConfigurationNodeConf()
	data := configuration_node_conf_model.List(c.Orm, page, c.PageSize, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// ID查询节点配置列表
func (c *ConfigNodeConfController) GetIdsToNodeConfList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	id := request_data["id"]

	configuration_node_conf_model := models.NewConfigurationNodeConf()
	data := configuration_node_conf_model.IdsToData(c.Orm, id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 节点配置删除
func (c *ConfigNodeConfController) DelNodeConf() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	node_conf_id, err := strconv.ParseInt(string(request_data["id"]), 10, 64)

	if (request_data["id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 验证是否有在使用,有在使用则禁止删除,未使用可以删除
	configuration_node_model := models.NewConfigurationNode()
	online_node := configuration_node_model.NodeConfIdOnLineCount(c.Orm, request_data["id"])
	if online_node[0]["count_ip"] != "0" {
		c.ResponseJson(helpers.CONFIGURATION_NODE_CONF_ONLINE_NODE_CODE, false, helpers.CONFIGURATION_NODE_CONF_ONLINE_NODE_MSG)
	}

	// TODO 删除不在线节点中的节点配置

	configuration_node_conf_model := models.NewConfigurationNodeConf()
	res, err := configuration_node_conf_model.Del(c.Orm, node_conf_id)

	if (err != nil) || (res == 0) {
		c.ResponseJson(helpers.CONFIGURATION_DEL_NODE_CONF_FAIL_CODE, false, helpers.CONFIGURATION_DEL_NODE_CONF_FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 节点配置修改
func (c *ConfigNodeConfController) EditNodeConf() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	id, err := strconv.ParseInt(string(request_data["id"]), 10, 64)
	conf := request_data["conf"]

	if (request_data["id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	configuration_node_conf_model := models.NewConfigurationNodeConf()
	configuration_node_conf_model.Id = int(id)
	data := map[string]interface{}{
		"conf": conf,
	}

	res := configuration_node_conf_model.Edit(c.Orm, data)

	if !res {
		c.ResponseJson(helpers.CONFIGURATION_EDIT_NODE_CONF_FAIL_CODE, false, helpers.CONFIGURATION_EDIT_NODE_CONF_FAIL_MSG)
	}

	// 获取使用该节点配置的ip地址
	configuration_node_model := models.NewConfigurationNode()
	ips := configuration_node_model.NodeConfIdToIp(c.Orm, request_data["id"])

	// 修改成功后需要下发数据
	go func() {
		var node_conf []configuration_constant.NodeConf
		json.Unmarshal([]byte(conf), &node_conf)

		// 向指定节点推送
		for _, val := range node_conf {
			for _, v := range val.Namespace {
				for i := 0; i < len(ips); i++ {
					node_conf_key := configuration_constant.CONFIGURATION_NODE_CONF + ips[i]["ip"].(string) + "/" + val.Appid + "_" + v.Name
					_, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), node_conf_key, v.Path)
				}
			}
		}

		// 节点初始化获取的数据修改
		for i := 0; i < len(ips); i++ {
			node_conf_key := configuration_constant.NODE_CONF_DIR + ips[i]["ip"].(string)
			_, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), node_conf_key, conf)
		}
	}()

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 节点配置添加
func (c *ConfigNodeConfController) AddNodeConf() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["name"] == "" || request_data["conf"] == "" {
		c.ResponseJson(helpers.CONFIGURATION_NODE_CONF_NAME_FAIL_CODE, false, helpers.CONFIGURATION_NODE_CONF_NAME_FAIL_MSG)
	}

	configuration_node_conf_model := models.NewConfigurationNodeConf()
	configuration_node_conf_model.AccountId = c.AccountInfo.Id
	configuration_node_conf_model.CreateTime = time.Now()
	configuration_node_conf_model.Name = request_data["name"]
	configuration_node_conf_model.Conf = request_data["conf"]

	//验证是否存在
	if len(configuration_node_conf_model.NameFind(c.Orm, request_data["name"])) > 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	node_conf_id, err := configuration_node_conf_model.Add(c.Orm)
	if err != nil {
		c.ResponseJson(helpers.CONFIGURATION_ADD_PROJECT_FAIL_CODE, false, helpers.CONFIGURATION_ADD_PROJECT_FAIL_MSG)
	}

	data := map[string]interface{}{
		"id":          node_conf_id,
		"username":    c.AccountInfo.Username,
		"name":        request_data["name"],
		"create_time": configuration_node_conf_model.CreateTime.Format("2006-01-02 15:04:05"),
		"nickname":    c.AccountInfo.Nickname,
		"conf":        request_data["conf"],
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}
