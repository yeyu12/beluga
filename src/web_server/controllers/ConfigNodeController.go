package controllers

import (
	"beluga/src/beluga/configuration_constant"
	"beluga/src/beluga/library"
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type ConfigNodeController struct {
	BaseController
}

// 节点列表
func (c *ConfigNodeController) GetNodeList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]

	if page == 0 {
		page = 1
	}

	configuration_node_model := models.NewConfigurationNode()
	data := configuration_node_model.List(c.Orm, page, c.PageSize, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)

}

// 删除节点
func (c *ConfigNodeController) DelNode() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	id, err := strconv.ParseInt(string(request_data["id"]), 10, 64)

	if (request_data["id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 只允许删除未在线的节点
	c.Orm.Begin()
	configuration_node_model := models.NewConfigurationNode()
	node_data := configuration_node_model.IdFind(c.Orm, id)
	res, err := configuration_node_model.Del(c.Orm, id)

	if (err != nil) || (res == 0) {
		c.Orm.Rollback()
		c.ResponseJson(helpers.CONFIGURATION_NODE_ONLINE_DEL_FATL_CODE, false, helpers.CONFIGURATION_NODE_ONLINE_DEL_FATL_MSG)
	}

	// 删除etcd中的数据
	etcd_key := configuration_constant.NODE_CONF_DIR + node_data.Ip
	_, err = library.G_conf_etcd_client.Kv.Delete(context.TODO(), etcd_key)

	if err != nil {
		c.Orm.Rollback()
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.Orm.Commit()

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 修改节点
func (c *ConfigNodeController) EditNode() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	id, err := strconv.ParseInt(string(request_data["id"]), 10, 64)
	node_conf_id := request_data["node_conf_id"]
	remake := request_data["remake"]

	if (request_data["id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.Orm.Begin()
	configuration_node_model := models.NewConfigurationNode()
	configuration_node_model.Id = id
	data := map[string]interface{}{
		"node_conf_id":           node_conf_id,
		"remake":                 remake,
		"conf_update_time":       time.Now(),
		"conf_update_account_id": c.AccountInfo.Id,
	}

	res := configuration_node_model.EditId(c.Orm, data)

	if !res {
		c.Orm.Rollback()
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 获取节点数据，判断是否在线，如果在线则下发通知，未在线则不通知
	// 写入etcd中, 编排为数组格式的key=val
	configuration_node_conf_model := models.NewConfigurationNodeConf()
	node_conf_data := configuration_node_conf_model.IdsToData(c.Orm, node_conf_id)
	var conf_data = make([]string, len(node_conf_data))

	for _, val := range node_conf_data {
		conf_data = append(conf_data, val["conf"].(string))
	}

	node_data := configuration_node_model.IdFind(c.Orm, id)
	if node_data.IsDelete != 0 {
		var conf_data []map[string]interface{}
		var node_conf_etcd []byte

		if node_data.NodeConfId != "" {
			configuration_node_conf_model := models.NewConfigurationNodeConf()
			node_conf_data := configuration_node_conf_model.IdsToData(orm.NewOrm(), node_data.NodeConfId)

			for _, val := range node_conf_data {
				var conf_map []map[string]interface{}
				json.Unmarshal([]byte(val["conf"].(string)), &conf_map)
				conf_data = append(conf_data, conf_map...)
			}

			node_conf_etcd, _ = json.Marshal(conf_data)
		}

		node_conf_key := configuration_constant.NODE_CONF_DIR + node_data.Ip
		_, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), node_conf_key, string(node_conf_etcd))

		if err != nil {
			beego.Error(err, "etcd节点配置发布失败")
		}
	}

	c.Orm.Commit()

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}
