package controllers

import (
	"beluga/src/web_server/helpers"
	"encoding/json"
	"strconv"
	"beluga/src/web_server/models"
	"time"
)

type NamespaceController struct {
	BaseController
}

func (c *NamespaceController) GetNamespaceList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]
	project_id, _ := strconv.Atoi(request_data["project_id"])

	if project_id == 0{
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	if page == 0 {
		page = 1
	}

	configuration_namespace_model := models.NewConfigurationNamespaceModel()
	data := configuration_namespace_model.List(c.Orm, page, c.PageSize, project_id, c.AccountInfo.Id, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 添加命名空间
func (c *NamespaceController) AddNamespace() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	if request_data["namespace_name"] == "" {
		c.ResponseJson(helpers.CONFIGURATION_NAMESPACE_NAME_FAIL_CODE, false, helpers.CONFIGURATION_NAMESPACE_NAME_FAIL_MSG)
	}
	if project_id == 0{
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	configuration_namespace_model := models.NewConfigurationNamespaceModel()

	// 去重
	namespace_data := configuration_namespace_model.GetProjectIdNamespaceNameToData(c.Orm, project_id, request_data["namespace_name"])
	if namespace_data.Id != 0 {
		c.ResponseJson(helpers.CONFIGURATION_NAMEPSACE_NAME_REPEAT_FAIL_CODE, false, helpers.CONFIGURATION_NAMEPSACE_NAME_REPEAT_FAIL_MSG)
	}

	configuration_namespace_model.AccountId = c.AccountInfo.Id
	configuration_namespace_model.CreateTime = time.Now()
	configuration_namespace_model.ProjectId = project_id
	configuration_namespace_model.NamespaceName = request_data["namespace_name"]

	namespace_id, err := configuration_namespace_model.AddNamespace(c.Orm)
	if err != nil {
		c.ResponseJson(helpers.CONFIGURATION_ADD_PROJECT_FAIL_CODE, false, helpers.CONFIGURATION_ADD_PROJECT_FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, map[string]int64{"id":namespace_id})
}

// 修改命名空间
func (c *NamespaceController) EditNamespace() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	namespace_id, _ := strconv.ParseInt(string(request_data["namespace_id"]), 10, 64)
	namespace_name := request_data["namespace_name"]

	if (request_data["namespace_id"] == "0") || (request_data["namespace_name"] == "") {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	configuration_namespace_model := models.NewConfigurationNamespaceModel()
	configuration_namespace_model.Id = int(namespace_id)
	data := map[string]interface{}{
		"namespace_name": namespace_name,
	}

	res := configuration_namespace_model.EditNamespace(c.Orm, data)

	if !res {
		c.ResponseJson(helpers.CONFIGURATION_NAMESPACE_EDIT_FAIL_CODE, false, helpers.CONFIGURATION_NAMESPACE_EDIT_FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 删除命名空间
func (c *NamespaceController) DelNamespace() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	namespace_id, err := strconv.ParseInt(string(request_data["namespace_id"]), 10, 64)

	if (request_data["namespace_id"] == "0") || (err != nil) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	configuration_namespace_model := models.NewConfigurationNamespaceModel()
	res, err := configuration_namespace_model.DelNamespace(c.Orm, namespace_id)

	if (err != nil) || (res == 0) {
		c.ResponseJson(helpers.CONFIGURATION_DEL_NAMESPACE_FAIL_CODE, false, helpers.CONFIGURATION_DEL_NAMESPACE_FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}
