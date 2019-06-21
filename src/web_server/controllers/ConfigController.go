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
	"github.com/coreos/etcd/clientv3"
	"strconv"
	"strings"
	"time"
)

type ConfigController struct {
	BaseController
}

// 获取详细配置(未发布的)
func (c *ConfigController) GetConfigLogList() {
	request_data := make(map[string]int)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id := int(request_data["project_id"])
	namespace_id := int(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}
	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)

	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_log_model := models.NewConfigurationLog()
	data := config_log_model.List(c.Orm, appid, namespace_name)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 获取已发布的配置
func (c *ConfigController) GetConfigList() {
	request_data := make(map[string]int)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id := int(request_data["project_id"])
	namespace_id := int(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}
	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)

	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_model := models.NewConfiguration()
	data := config_model.List(c.Orm, appid, namespace_name)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 添加配置
func (c *ConfigController) AddConfig() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])
	config_key := string(request_data["key"])
	config_val := string(request_data["val"])
	config_remake := string(request_data["remake"])

	if (project_id == 0) || (namespace_id == 0) || (config_key == "") || (config_val == "") {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.Orm.Begin()
	res_sta, res_str := c.setAddConfigLog(appid, namespace_name, config_key, config_val, config_remake)
	if res_sta {
		c.Orm.Commit()
	}

	c.ResponseStr(res_str)
}

// 修改配置
func (c *ConfigController) EditConfig() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])
	is_release, _ := strconv.Atoi(request_data["release_status"])
	config_key := string(request_data["key"])
	config_val := string(request_data["val"])
	config_remake := string(request_data["remake"])

	if (project_id == 0) || (namespace_id == 0) || (config_key == "") || (config_val == "") {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_log_model := models.NewConfigurationLog()
	configuration_version_model, version_data := c.getVersion(appid, namespace_name)
	version := version_data.Version

	if is_release != 0 {
		if version == "" {
			version = c.createVersion(configuration_version_model, appid, namespace_name).Version
		} else {
			// 查找未发布的配置
			config_log_data := config_log_model.KeyFind(c.Orm, appid, namespace_name, version, config_key, 0)
			if config_log_data.Key != "" {
				c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
			}
		}

		config_model := models.NewConfiguration()
		config_release_data := config_model.KeyFind(c.Orm, appid, namespace_name, config_key)

		if config_release_data.Key != "" {
			if config_remake == "" {
				config_remake = config_release_data.Remake
			}
			// 构建数据
			config_log_model.Id = 0
			config_log_model.Appid = appid
			config_log_model.NamespaceName = namespace_name
			config_log_model.Version = version
			config_log_model.AccountId = c.AccountInfo.Id
			config_log_model.Key = config_key
			config_log_model.Val = config_val
			config_log_model.Remake = config_remake
			config_log_model.Type = 2
			config_log_model.CreateTime = time.Now()
			config_log_model.UpdateTime = time.Now()
			config_log_model.OldVal = config_release_data.Val
			config_log_model.OldRemake = config_release_data.Remake

			// 添加数据
			if _, err := config_log_model.Add(c.Orm); err != nil {
				c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
			}
		} else {
			c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
		}
	} else {
		config_log_data := config_log_model.KeyFind(c.Orm, appid, namespace_name, version, config_key, 0)

		if config_log_data.Key != "" {
			config_log_model.Id = config_log_data.Id
			data := map[string]interface{}{
				"val":         config_val,
				"remake":      config_remake,
				"update_time": time.Now(),
			}

			if !config_log_model.Edit(c.Orm, data) {
				c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
			}
		}
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 删除配置
func (c *ConfigController) DelConfig() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])
	is_release, _ := strconv.Atoi(request_data["is_release"])
	config_key := string(request_data["key"])

	if (project_id == 0) || (namespace_id == 0) || (config_key == "") {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	_,res_str := c.setDelConfigLog(appid, namespace_name, config_key, is_release)
	c.ResponseStr(res_str)
}

// 发布
func (c *ConfigController) ReleaseConfig() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.Orm.Begin()
	res_sta, res_str := c.setReleaseConfig(appid, namespace_name, 0)
	if res_sta {
		c.Orm.Commit()

		// 同步配置
		go c.syncEtcd(appid, namespace_name)
	}

	c.ResponseStr(res_str)
}

// 获取回滚数据
func (c *ConfigController) GetRollbackLast() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 获取到上次发布之前的版本号，回滚的不算
	config_operation_model := models.NewConfigurationOperation()
	operation_data := config_operation_model.LastRecords(c.Orm, appid, namespace_name)
	config_log_mode := models.NewConfigurationLog()
	data := config_log_mode.VersionToConfigToReleaseLog(c.Orm, appid, namespace_name, operation_data.Version)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 回滚
func (c *ConfigController) RollbackConfig() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 获取到上次发布之前的版本号，回滚的不算
	config_operation_model := models.NewConfigurationOperation()
	operation_data := config_operation_model.LastRecords(c.Orm, appid, namespace_name)
	config_log_mode := models.NewConfigurationLog()
	config_log_release_data := config_log_mode.VersionToConfigToReleaseLog(c.Orm, appid, namespace_name, operation_data.Version)

	if len(config_log_release_data) == 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 获取当前未发布的配置数据
	config_version_model := models.NewConfigurationVersion()
	config_version_data := config_version_model.GetVersion(c.Orm, appid, namespace_name)

	// 当前未发布的本版操作的数据
	var config_log_data []orm.Params
	c.Orm.Begin()

	if config_version_data.Version == "" {
		config_version_data = *c.createVersion(config_version_model, appid, namespace_name)
		goto JAMP
	}

	config_log_data = config_log_mode.VersionToConfigLog(c.Orm, appid, namespace_name, config_version_data.Version)
	// 清空数据
	if _, err := config_log_mode.DelLog(c.Orm, appid, namespace_name, config_version_data.Version); err != nil {
		c.Orm.Rollback()
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

JAMP:

// 上个版本已发布的数据，写入到当前版本中去
	for _, val := range config_log_release_data {
		config_log_release_type, _ := strconv.Atoi(val["type"].(string))

		config_log_mode.Id = 0
		config_log_mode.Appid = appid
		config_log_mode.NamespaceName = namespace_name
		config_log_mode.Version = config_version_data.Version
		config_log_mode.AccountId = c.AccountInfo.Id
		config_log_mode.IsRelease = 0
		config_log_mode.CreateTime = time.Now()
		config_log_mode.UpdateTime = time.Now()

		switch config_log_release_type {
		case -1:
			config_log_mode.Key = val["key"].(string)
			config_log_mode.Val = val["val"].(string)
			config_log_mode.Remake = val["remake"].(string)
			config_log_mode.OldVal = val["old_val"].(string)
			config_log_mode.OldRemake = val["old_remake"].(string)
			config_log_mode.Type = 2

			config_log_mode.Add(c.Orm)
			break
		case 1:
			config_log_mode.Key = val["key"].(string)
			config_log_mode.Val = val["val"].(string)
			config_log_mode.Remake = val["remake"].(string)
			config_log_mode.OldVal = val["old_val"].(string)
			config_log_mode.OldRemake = val["old_remake"].(string)
			config_log_mode.Type = -1

			config_log_mode.Add(c.Orm)
			break
		case 2:
			config_log_mode.Key = val["key"].(string)
			config_log_mode.Val = val["old_val"].(string)
			config_log_mode.Remake = val["old_remake"].(string)
			config_log_mode.OldVal = ""
			config_log_mode.OldRemake = ""
			config_log_mode.Type = 2

			config_log_mode.Add(c.Orm)
			break
		}
	}

	// 发布版本
	res_sta, res_str := c.setReleaseConfig(appid, namespace_name, -1)
	if !res_sta {
		c.Orm.Rollback()
		c.ResponseStr(res_str)
	}

	// 已发布的操作修改为回滚
	config_operation_model.EditOperationRollback(c.Orm, appid, namespace_name, operation_data.Version)

	// 合并数据，数据写入

	// 获取当前未发布的版本号
	configuration_version_model, version_data := c.getVersion(appid, namespace_name)
	version := version_data.Version

	if version == "" {
		// 创建新的版本号
		version = c.createVersion(configuration_version_model, appid, namespace_name).Version
	}

	// 上个版本已发布的数据写入到当前未发布的数据中去
	for _, val := range config_log_release_data {
		config_log_mode.Id = 0
		config_log_mode.Appid = appid
		config_log_mode.NamespaceName = namespace_name
		config_log_mode.Version = version
		config_log_mode.AccountId = c.AccountInfo.Id
		config_log_mode.IsRelease = 0
		config_log_mode.CreateTime = time.Now()
		config_log_mode.UpdateTime = time.Now()
		config_log_mode.Key = val["key"].(string)
		config_log_mode.Val = val["val"].(string)
		config_log_mode.Remake = val["remake"].(string)
		config_log_mode.OldVal = val["old_val"].(string)
		config_log_mode.OldRemake = val["old_remake"].(string)
		config_log_mode.Type, _ = strconv.Atoi(val["type"].(string))

		config_log_mode.Add(c.Orm)
	}

	// 当前版本未发布的数据,写入到新版本中去
	for _, val := range config_log_data {
		config_log_type, _ := strconv.Atoi(val["type"].(string))

		config_log_mode.Id = 0
		config_log_mode.Appid = appid
		config_log_mode.NamespaceName = namespace_name
		config_log_mode.Version = version
		config_log_mode.AccountId = c.AccountInfo.Id
		config_log_mode.IsRelease = 0
		config_log_mode.Type = config_log_type
		config_log_mode.CreateTime = time.Now()
		config_log_mode.UpdateTime = time.Now()

		// 检测当前配置是否存在
		config_log := config_log_mode.KeyFind(c.Orm, appid, namespace_name, version, val["key"].(string), 0)

		// 合并规则
		// 添加，未操作
		// 添加，删除-删除
		// 添加，修改
		// 修改，未操作
		// 修改，修改
		// 修改，删除-删除
		// 删除，未操作
		// 删除，新增-删除
		switch config_log_type {
		case -1:
			if config_log.Key != "" {
				config_log_mode.AppidNamesapceNameKeyToEdit(c.Orm, appid, namespace_name, val["key"].(string), map[string]interface{}{
					"type": -1,
				})
				goto CONFIG_ADD_FAIL_JAMP
			}

			break
		case 1:
			if config_log.Key != "" {
				switch config_log.Type {
				case -1:
					config_log_mode.AppidNamesapceNameKeyToEdit(c.Orm, appid, namespace_name, val["key"].(string), map[string]interface{}{
						"type": -1,
					})
					break
				case 2:
					config_log_mode.AppidNamesapceNameKeyToEdit(c.Orm, appid, namespace_name, val["key"].(string), map[string]interface{}{
						"type":   1,
						"val":    val["val"].(string),
						"remake": val["remake"].(string),
					})

					break
				}

				goto CONFIG_ADD_FAIL_JAMP
			}

			break
		case 2:
			if config_log.Key != "" {
				switch config_log.Type {
				case -1:
					config_log_mode.AppidNamesapceNameKeyToEdit(c.Orm, appid, namespace_name, val["key"].(string), map[string]interface{}{
						"type": -1,
					})
					break
				case 2:
					config_log_mode.AppidNamesapceNameKeyToEdit(c.Orm, appid, namespace_name, val["key"].(string), map[string]interface{}{
						"type":   2,
						"val":    val["val"].(string),
						"remake": val["remake"].(string),
					})
					break
				}

				goto CONFIG_ADD_FAIL_JAMP
			}

			break
		}

		config_log_mode.Key = val["key"].(string)
		config_log_mode.Val = val["val"].(string)
		config_log_mode.Remake = val["remake"].(string)
		config_log_mode.OldVal = val["old_val"].(string)
		config_log_mode.OldRemake = val["old_remake"].(string)
		config_log_mode.Add(c.Orm)

	CONFIG_ADD_FAIL_JAMP:
	}

	// 事物提交完成。
	c.Orm.Commit()
	// 同步配置
	go c.syncEtcd(appid, namespace_name)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, config_log_release_data)
}

// 发布历史
func (c *ConfigController) ReleaseHistory() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_opreation_model := models.NewConfigurationOperation()
	config_operation := config_opreation_model.List(c.Orm, appid, namespace_name)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, config_operation)
}

// 获取项目名、命名空间名
func (c *ConfigController) GetProjectNameToNamespaceName() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	// 获取项目名和命名空间名称
	config_project_mode := models.NewConfigurationProjectModel()
	project_data := config_project_mode.GetProjectIdToData(c.Orm, project_id)
	config_namespace_model := models.NewConfigurationNamespaceModel()
	namespace_data := config_namespace_model.GetNamespaceIdToData(c.Orm, project_id, namespace_id)

	data := map[string]interface{}{
		"project_name":   project_data.ProjectName,
		"namespace_name": namespace_data.NamespaceName,
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 获取变更的配置
func (c *ConfigController) GetConfigVersionReleaseList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])
	version := request_data["version"]

	if (project_id == 0) || (namespace_id == 0) || (version == "") {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_log_model := models.NewConfigurationLog()
	data := config_log_model.VersionToConfigToReleaseLog(c.Orm, appid, namespace_name, version)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 获取版本号下所有的配置数据
func (c *ConfigController) GetConfigAllVersionReleaseList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])
	version := request_data["version"]

	if (project_id == 0) || (namespace_id == 0) || (version == "") {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_operation_log_model := models.NewConfigurationOperationLog()
	data := config_operation_log_model.VersionToConfigToReleaseLog(c.Orm, appid, namespace_name, version)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 配置文本提交
func (c *ConfigController) SaveConfigText() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])
	config_text_data := request_data["config_data"]

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_text_arr := strings.Split(config_text_data, "\n")
	config_text_map := make(map[string]interface{})

	// 数据写入
	for _, val := range config_text_arr {
		config_temp := strings.Split(val, "=")
		config_temp_len := len(config_temp)

		if (config_temp_len == 0) || (config_temp_len == 1) {
			continue
		}

		config_temp[0] = strings.Trim(config_temp[0], " ")
		config_temp[1] = strings.Trim(config_temp[1], " ")
		config_text_map[config_temp[0]] = config_temp[1]

		c.setAddConfigLog(appid, namespace_name, config_temp[0], config_temp[1], "")
	}

	// 处理删除的数据
	// 获取已发布的数据
	config_model := models.NewConfiguration()
	config_data := config_model.List(c.Orm, appid, namespace_name)

	for _, val := range config_data {
		if config_text_map[val["key"].(string)] == nil {
			c.setDelConfigLog(appid, namespace_name, val["key"].(string), 1)
		}
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, config_text_arr)
}

// 配置同步
// 同步之后，需要下发同步的配置信息
func (c *ConfigController) SyncConfig()  {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	project_id, _ := strconv.Atoi(request_data["project_id"])
	namespace_id, _ := strconv.Atoi(request_data["namespace_id"])

	if (project_id == 0) || (namespace_id == 0) {
		c.ResponseJson(helpers.FAIL_PARAMS_CODE, false, helpers.FAIL_PARAMS_MSG)
	}

	appid, namespace_name := c.projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id)
	if appid == "" || namespace_name == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	go c.syncEtcd(appid, namespace_name)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 获取appid和namespace_name
func (c *ConfigController) projectIdNamespaceIdToAppidNamespaceName(project_id, namespace_id int) (string, string) {
	configuration_project_model := models.NewConfigurationProjectModel()
	appid := configuration_project_model.GetProjectIdToAppid(c.Orm, project_id)

	configuration_namespace_mode := models.NewConfigurationNamespaceModel()
	namespace_name := configuration_namespace_mode.GetNamespaceIdToName(c.Orm, project_id, namespace_id)

	return appid, namespace_name
}

// 生成新版本
func (c *ConfigController) createVersion(configuration_version_model *models.ConfigurationVersion, appid, namespace_name string) *models.ConfigurationVersion {
	config_release_data := configuration_version_model.GetReleaseVersion(c.Orm, appid, namespace_name)
	if config_release_data.Id == 0 {
		configuration_version_model.PId = 0
	} else {
		configuration_version_model.PId = config_release_data.Id
	}

	configuration_version_model.Appid = appid
	configuration_version_model.NamespaceName = namespace_name
	configuration_version_model.Version = helpers.GenVersion()
	configuration_version_model.CreateTime = time.Now()

	if _, err := configuration_version_model.Add(c.Orm); err != nil {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	return configuration_version_model
}

// appid和namespace_name获取最新未发布的版本数据
func (c *ConfigController) getVersion(appid, namespace_name string) (*models.ConfigurationVersion, models.ConfigurationVersion) {
	configuration_version_model := models.NewConfigurationVersion()

	return configuration_version_model, configuration_version_model.GetVersion(c.Orm, appid, namespace_name)
}

// 发布方法
func (c *ConfigController) setReleaseConfig(appid, namespace_name string, operation_type int) (bool, string) {
	config_version_model, version_data := c.getVersion(appid, namespace_name)

	if version_data.Version == "" {
		return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	config_log_model := models.NewConfigurationLog()
	config_model := models.NewConfiguration()

	// 查找出所有数据
	config_log_data := config_log_model.VersionToConfigLog(c.Orm, appid, namespace_name, version_data.Version)

	config_model.Appid = appid
	config_model.NamespaceName = namespace_name

	config_operation_model := models.NewConfigurationOperation()
	config_operation_log_model := models.NewConfigurationOperationLog()

	config_operation_model.Appid = appid
	config_operation_model.NamespaceName = namespace_name
	config_operation_model.Version = version_data.Version
	config_operation_model.AccountId = c.AccountInfo.Id
	config_operation_model.OperationType = operation_type
	config_operation_model.Rollback = operation_type
	config_operation_model.CreateTime = time.Now()

	config_operation_id, err := config_operation_model.Add(c.Orm)
	config_operation_log_model.ConfigurationOperationId = int(config_operation_id)

	if err != nil {
		beego.Error("配置发布失败!操作记录添加失败")
		c.Orm.Rollback()
		return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	for _, val := range config_log_data {
		release_type, _ := strconv.Atoi(val["type"].(string))

		switch release_type {
		case -1:
			config_model.Del(c.Orm, appid, namespace_name, val["key"].(string))
			break
		case 1:
			config_model.Id = 0
			config_model.Key = val["key"].(string)
			config_model.Val = val["val"].(string)
			config_model.Remake = val["remake"].(string)
			config_model.CreateTime = time.Now()

			config_model.Add(c.Orm)
			break
		case 2:
			config_model.AppidNamesapceNameKeyToEdit(c.Orm, appid, namespace_name, val["key"].(string), map[string]interface{}{
				"val":    val["val"],
				"remake": val["remake"],
			})
			break
		}
	}

	// 获取所有已发布的配置,写入
	for _, val := range config_model.List(c.Orm, appid, namespace_name) {
		config_operation_log_model.Id = 0
		config_operation_log_model.AccountId = c.AccountInfo.Id
		config_operation_log_model.Key = val["key"].(string)
		config_operation_log_model.Val = val["val"].(string)
		config_operation_log_model.Remake = val["remake"].(string)
		config_operation_log_model.CreateTime = time.Now()

		config_operation_log_model.Add(c.Orm)
	}

	// 更新纪录状态
	if !config_log_model.AppidNamespaceNameVersionToEdit(c.Orm, appid, namespace_name, map[string]interface{}{"is_release": 1}) {
		beego.Error("配置发布失败！配置操作记录失败")
		c.Orm.Rollback()
		return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 更新版本号状态
	config_version_model.Id = version_data.Id
	if !config_version_model.Update(c.Orm, map[string]interface{}{"is_release": 1}) {
		beego.Error("配置发布失败！修改版本状态失败")
		c.Orm.Rollback()
		return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	go c.syncEtcd(appid, namespace_name)

	return true, c.ResStr(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 添加配置方法
func (c *ConfigController) setAddConfigLog(appid, namespace_name, config_key, config_val, config_remake string) (bool, string) {
	configuration_version_model, version_data := c.getVersion(appid, namespace_name)
	version := version_data.Version

	if version == "" {
		// 创建新的版本号
		version = c.createVersion(configuration_version_model, appid, namespace_name).Version
	}

	config_log_model := models.NewConfigurationLog()

	// 查找未发布的配置
	config_log_data := config_log_model.KeyFind(c.Orm, appid, namespace_name, version, config_key, 0)

	// 在未发布里面
	if config_log_data.Key != "" {
		config_log_model = &config_log_data
		config_log_edit_data := map[string]interface{}{
			"val":         config_val,
			"update_time": time.Now(),
			"account_id":  c.AccountInfo.Id,
		}

		if !config_log_model.Edit(c.Orm, config_log_edit_data) {
			c.Orm.Rollback()
			return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
		}
	} else {
		// 查找发布了的配置
		config_model := models.NewConfiguration()
		config_release_data := config_model.KeyFind(c.Orm, appid, namespace_name, config_key)

		if config_release_data.Key != "" {
			// 如果值相同则跳出不添加
			if config_release_data.Val == config_val {
				goto EXISTENCE_JUMP
			}
			config_log_model.Type = 2
			config_log_model.Key = config_release_data.Key
			config_log_model.OldVal = config_release_data.Val
			config_log_model.OldRemake = config_release_data.Remake

			goto JUMP
		}
		// 添加配置记录,构建数据
		config_log_model.Type = 1
		config_log_model.Key = config_key
	JUMP:
		config_log_model.Version = version
		config_log_model.Appid = appid
		config_log_model.NamespaceName = namespace_name
		config_log_model.AccountId = c.AccountInfo.Id
		config_log_model.Val = config_val
		config_log_model.Remake = config_remake
		config_log_model.CreateTime = time.Now()
		config_log_model.UpdateTime = config_log_model.CreateTime

		// 添加数据
		if _, err := config_log_model.Add(c.Orm); err != nil {
			c.Orm.Rollback()
			return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
		}
	}

EXISTENCE_JUMP:

	return true, c.ResStr(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 删除配置方法
func (c *ConfigController) setDelConfigLog(appid, namespace_name, config_key string, is_release int) (bool, string) {
	config_log_model := models.NewConfigurationLog()
	configuration_version_model, version_data := c.getVersion(appid, namespace_name)
	version := version_data.Version

	// 配置是否发布
	if is_release != 0 {
		// 版本不存在添加版本号
		if version == "" {
			version = c.createVersion(configuration_version_model, appid, namespace_name).Version
		} else {
			// 查找未发布的配置
			config_log_data := config_log_model.KeyFind(c.Orm, appid, namespace_name, version, config_key, 0)
			if config_log_data.Key != "" {
				return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
			}
		}

		config_model := models.NewConfiguration()
		config_release_data := config_model.KeyFind(c.Orm, appid, namespace_name, config_key)

		if config_release_data.Key != "" {
			// 构建数据
			config_log_model.Id = 0
			config_log_model.Appid = appid
			config_log_model.NamespaceName = namespace_name
			config_log_model.Version = version
			config_log_model.AccountId = c.AccountInfo.Id
			config_log_model.Key = config_key
			config_log_model.Val = config_release_data.Val
			config_log_model.Type = -1
			config_log_model.CreateTime = time.Now()
			config_log_model.UpdateTime = time.Now()

			// 添加数据
			if _, err := config_log_model.Add(c.Orm); err != nil {
				return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
			}
		} else {
			return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
		}
	} else {
		// 查找未发布的配置
		config_log_data := config_log_model.KeyFind(c.Orm, appid, namespace_name, version, config_key, 0)

		if config_log_data.Key != "" {
			// 在未发布的数据中存在,删除该记录
			if _, err := config_log_model.Del(c.Orm, appid, namespace_name, config_key); err != nil {
				return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
			}
		} else {
			return false, c.ResStr(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
		}
	}

	return true, c.ResStr(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 同步到etcd中
func (c *ConfigController)syncEtcd(appid, namespace string)  {
	lease_grant_resp, err := library.G_conf_etcd_client.Lease.Grant(context.TODO(), 86400)
	if err != nil {
		beego.Error("配置同步失败，创建租约失败。")
		return
	}

	key := configuration_constant.CONFIGURATION_CONF_SYNC + appid + "_" + namespace
	_, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), key, "", clientv3.WithLease(lease_grant_resp.ID))

	if err != nil {
		beego.Error("配置同步失败")
	}
}