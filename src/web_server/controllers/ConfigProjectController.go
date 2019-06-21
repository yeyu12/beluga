package controllers

import (
	"encoding/json"
	"beluga/src/web_server/models"
	"beluga/src/web_server/helpers"
	"time"
	"strconv"
)

type ConfigProjectController struct {
	BaseController
}

// 获取项目列表
func (c *ConfigProjectController) GetProjectList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]

	if page == 0 {
		page = 1
	}

	configuration_project_model := models.NewConfigurationProjectModel()
	data := configuration_project_model.List(c.Orm, page, c.PageSize, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 获取用户自有项目列表
func (c *ConfigProjectController) GetUserProjectList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])

	if page == 0 {
		page = 1
	}

	configuration_project_model := models.NewConfigurationProjectModel()
	data := configuration_project_model.UserProjectList(c.Orm, page, c.PageSize, c.AccountInfo.Id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 添加项目
func (c *ConfigProjectController) AddProject() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["project_name"] == "" {
		c.ResponseJson(helpers.CONFIGURATION_PROJECT_NAME_FAIL_CODE, false, helpers.CONFIGURATION_PROJECT_NAME_FAIL_MSG)
	}

	configuration_project_model := models.NewConfigurationProjectModel()

	// TODO 考虑项目名称到底要不要去重，因为appid不同
	project_data := configuration_project_model.GetProjectNameToData(c.Orm, request_data["project_name"])
	if project_data.Id != 0 {
		c.ResponseJson(helpers.CONFIGURATION_PROJECT_NAME_REPEAT_FAIL_CODE, false, helpers.CONFIGURATION_PROJECT_NAME_REPEAT_FAIL_MSG)
	}

	configuration_project_model.AccountId = c.AccountInfo.Id
	configuration_project_model.CreateTime = time.Now()
	configuration_project_model.ProjectName = request_data["project_name"]
	configuration_project_model.Appid = helpers.RandStr()

	c.Orm.Begin()

	project_id, err := configuration_project_model.AddProject(c.Orm)
	if err != nil {
		c.Orm.Rollback()
		c.ResponseJson(helpers.CONFIGURATION_ADD_PROJECT_FAIL_CODE, false, helpers.CONFIGURATION_ADD_PROJECT_FAIL_MSG)
	}

	// 项目数+1
	account_model := models.NewAccount()
	c.AccountInfo.ConfigurationNum++
	if !account_model.Update(c.Orm, c.AccountInfo) {
		c.Orm.Rollback()
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	configuration_project_model.Id = int(project_id)

	c.Orm.Commit()
	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, configuration_project_model)
}

// 删除项目
func (c *ConfigProjectController) DelProject() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, err := strconv.ParseInt(string(request_data["project_id"]), 10, 64)

	if (request_data["project_id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.Orm.Begin()
	configuration_project_model := models.NewConfigurationProjectModel()
	res, err := configuration_project_model.DelProject(c.Orm, project_id)

	if (err != nil) || (res == 0) {
		c.Orm.Rollback()
		c.ResponseJson(helpers.CONFIGURATION_DEL_PROJECT_FAIL_CODE, false, helpers.CONFIGURATION_DEL_PROJECT_FAIL_MSG)
	}

	// 项目数+1
	account_model := models.NewAccount()
	c.AccountInfo.ConfigurationNum--
	if !account_model.Update(c.Orm, c.AccountInfo) {
		c.Orm.Rollback()
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.Orm.Commit()
	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 修改项目
func (c *ConfigProjectController) EditProject() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, err := strconv.ParseInt(string(request_data["project_id"]), 10, 64)
	project_name := request_data["project_name"]

	if (request_data["project_id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	configuration_project_model := models.NewConfigurationProjectModel()
	configuration_project_model.Id = int(project_id)
	data := map[string]interface{}{
		"project_name": project_name,
	}

	res := configuration_project_model.EditProject(c.Orm, data)

	if !res {
		c.ResponseJson(helpers.CONFIGURATION_EDIT_PROJECT_FAIL_CODE, false, helpers.CONFIGURATION_EDIT_PROJECT_FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}
