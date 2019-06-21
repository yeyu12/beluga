package controllers

import (
	"beluga/src/web_server/helpers"
	"encoding/json"
	"beluga/src/web_server/models"
	"strconv"
	"time"
	"crypto/md5"
	"encoding/hex"
)

type AccountController struct {
	BaseController
}

// 获取个人信息
func (c *AccountController) GetUserInfo() {
	account_info := c.AccountInfo
	account_info.Passwd = ""
	account_info.Token = ""

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, account_info)
}

// 设置个人信息
func (c *AccountController) SetUserInfo() {
	request_data := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	username := request_data["username"].(string)
	nickname := request_data["nickname"].(string)

	if username == "" || nickname == "" {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	account_model := models.NewAccount()
	account_model = c.AccountInfo
	account_model.Username = username
	account_model.Nickname = nickname

	if !account_model.Update(c.Orm, account_model) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	account_model.Token = ""
	account_model.Passwd = ""

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, account_model)
}

// 修改密码
func (c *AccountController) ChangePasswd()  {
	request_data := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	passwd := request_data["passwd"].(string)
	new_passwd := request_data["new_passwd"].(string)

	if passwd == "" || new_passwd == "" {
		c.ResponseJson(helpers.CHANGE_PASSWD_EMPTY_FAIL_CODE, false, helpers.CHANGE_PASSWD_EMPTY_FAIL_MSG)
	}

	if passwd == new_passwd {
		c.ResponseJson(helpers.CHANGE_PASSWD_IDENTICAL_FAIL_CODE, false, helpers.CHANGE_PASSWD_IDENTICAL_FAIL_MSG)
	}

	if c.AccountInfo.Passwd != helpers.EncryptionPasswd(passwd) {
		c.ResponseJson(helpers.CHANGE_PASSWD_FAIL_CODE, false, helpers.CHANGE_PASSWD_FAIL_MSG)
	}

	if c.AccountInfo.Passwd == helpers.EncryptionPasswd(new_passwd) {
		c.ResponseJson(helpers.CHANGE_PASSWD_IDENTICAL_FAIL_CODE, false, helpers.CHANGE_PASSWD_IDENTICAL_FAIL_MSG)
	}

	account_model := models.NewAccount()
	account_model = c.AccountInfo
	account_model.Passwd = helpers.EncryptionPasswd(new_passwd)

	token_str := c.AccountInfo.Username + strconv.Itoa(int(time.Now().Unix()))

	ctx := md5.New()
	ctx.Write([]byte(token_str))
	token := hex.EncodeToString(ctx.Sum(nil))

	account_model.Token = token
	account_model.TokenExpiryTime = time.Now()

	if !account_model.Update(c.Orm, account_model) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}
