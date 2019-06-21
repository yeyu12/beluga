package controllers

import (
	"encoding/json"
	"strconv"
	"beluga/src/web_server/models"
	"beluga/src/web_server/helpers"
	"time"
)

type UserController struct {
	BaseController
}

// 用户列表
func (c *UserController) List()  {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]

	if page == 0 {
		page = 1
	}

	account_model := models.NewAccount()
	data := account_model.List(c.Orm, page, c.PageSize, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 添加用户
func (c *UserController)Add()  {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["username"] == "" || request_data["passwd"] == "" {
		c.ResponseJson(helpers.USERNAME_TO_PASSWD_FAIL_CODE, false, helpers.USERNAME_TO_PASSWD_FAIL_MSG)
	}

	account_model := models.NewAccount()

	// 判断用户是否存在
	account_data, _ := account_model.UsernameToData(c.Orm, request_data["username"])
	if account_data.Id != 0 {
		c.ResponseJson(helpers.USERNAME_EXISTENCE_CODE, false, helpers.USERNAME_EXISTENCE_MSG)
	}

	account_model.Username = request_data["username"]
	account_model.Passwd = helpers.EncryptionPasswd(request_data["passwd"])
	account_model.CreateTime = time.Now()
	account_model.Nickname = account_model.Username
	account_model.Status = 1
	account_id, err := account_model.Add(c.Orm)

	if err != nil {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	account_model.Id = int(account_id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, account_model)
}

// 修改用户信息
func (c *UserController) Edit()  {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["username"] == ""{
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	username := request_data["username"]
	data := make(map[string]interface{})

	data["status"] = request_data["status"]
	if request_data["nickname"] != "" {
		data["nickname"] = string(request_data["nickname"])
	}
	if request_data["passwd"] != "" {
		data["passwd"] = helpers.EncryptionPasswd(request_data["passwd"])
	}

	account_model := models.NewAccount()
	if !account_model.IdSetUserInfo(c.Orm, username, data) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}