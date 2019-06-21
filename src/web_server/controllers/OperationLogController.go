package controllers

import (
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"encoding/json"
	"strconv"
)

type OperationLogController struct {
	BaseController
}

// 操作记录
func (c *OperationLogController) List()  {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	account_id := 0
	page, _ := strconv.Atoi(request_data["page"])
	nickname := request_data["search"]

	if page == 0 {
		page = 1
	}
	if nickname != "" {
		account_model := models.NewAccount()
		account_info, _ := account_model.UserNicknameToInfo(c.Orm, nickname)
		account_id = account_info.Id
	}

	configuration_project_model := models.NewOperationLog()
	data := configuration_project_model.List(c.Orm, page, c.PageSize, int64(account_id))

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}