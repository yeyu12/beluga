package controllers

import (
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"encoding/json"
	"strconv"
)

type TaskLogController struct {
	BaseController
}

// 列表
func (c *BaseController) List() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	task_id := request_data["task_id"]

	if page == 0 {
		page = 1
	}

	task_model := models.NewTaskLog()
	data := task_model.List(c.Orm, page, c.PageSize, task_id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}
