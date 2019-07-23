package controllers

import (
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"encoding/json"
	"strconv"
)

type TaskNodeController struct {
	BaseController
}

// 节点列表
func (c *TaskNodeController) GetNodeList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]

	if page == 0 {
		page = 1
	}

	task_node_model := models.NewTaskNode()
	data := task_node_model.List(c.Orm, page, c.PageSize, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 删除节点
func (c *TaskNodeController) Del()  {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	id, err := strconv.ParseInt(string(request_data["id"]), 10, 64)

	if (request_data["id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 只允许删除未在线的节点
	task_node_model := models.NewTaskNode()
	res, err := task_node_model.Del(c.Orm, id)

	if (err != nil) || (res == 0) {
		c.ResponseJson(helpers.TASK_NODE_ONLINE_DEL_FATL_CODE, false, helpers.TASK_NODE_ONLINE_DEL_FATL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}