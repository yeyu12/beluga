package controllers

import (
	"github.com/mojocn/base64Captcha"
	"github.com/astaxie/beego"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"time"
	"strconv"
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
)

type LoginController struct {
	BaseController
}

// 登录结构体
type login struct {
	Username          string `json:"username"`
	Passwd            string `json:"passwd"`
	Verification_code string `json:"verification_code"`
	Remember_passwd   bool   `json:"remember_passwd"`
	Captcha_id        string `json:"captcha_id"`
}

// 登录
func (c *LoginController) Login() {
	if c.Ctx.Input.IsAjax() {
		if c.AccountInfo != nil {
			c.ResponseJson(helpers.REPEAT_LOGIN_FAIL_CODE, false, helpers.REPEAT_LOGIN_FAIL_MSG)
		}

		var login_data map[string]string
		json.Unmarshal(c.Ctx.Input.RequestBody, &login_data)

		username := login_data["username"]
		passwd := helpers.EncryptionPasswd(login_data["passwd"])

		// 验证码验证
		//verification_code := login_data["Verification_code"]
		//captcha_id := login_data["captcha_id"]

		/*if captcha_id == "" {
			c.ResponseJson(helpers.ILLEGAL_FAIL_CODE, false, helpers.ILLEGAL_FAIL_MSG)
		}

		if !c.VerfiyCaptcha(captcha_id, verification_code) {
			c.ResponseJson(helpers.CAPTCHA_FAIL_CODE, false, helpers.CAPTCHA_FAIL_MSG)
		}*/

		account_model := models.NewAccount()
		account_data, e := account_model.UsernameFind(c.Orm, username, passwd)

		if e != nil {
			c.ResponseJson(helpers.USERNAME_PASSWD_ERROR_CODE, false, helpers.USERNAME_PASSWD_ERROR_MSG)
		}

		if account_data.Status != 1 {
			c.ResponseJson(helpers.USER_PROHIBIT_CODE, false, helpers.USER_PROHIBIT_MSG)
		}

		// 修改token，并更新token过期时间
		token_str := username + strconv.Itoa(int(time.Now().Unix()))

		ctx := md5.New()
		ctx.Write([]byte(token_str))
		token := hex.EncodeToString(ctx.Sum(nil))

		account_data.Token = token
		account_data.TokenExpiryTime = time.Now()

		if !account_model.Update(c.Orm, account_data) {
			c.ResponseJson(helpers.LOGIN_TOKEN_UPDATE_FATL_CODE, false, helpers.LOGIN_TOKEN_UPDATE_FATL_MSG)
		}

		c.AccountInfo = account_data

		c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, map[string]string{"token": account_data.Token})
	} else {
		c.ResponseJson(helpers.ILLEGAL_FAIL_CODE, false, helpers.ILLEGAL_FAIL_MSG)
	}
}

// 退出登录
func (c *LoginController) Logout() {
	c.DelSession("user_info")

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 忘记密码
func (c *LoginController) ForgotPassword() {

}

// 生成图片验证码
func (c *LoginController) Captcha() {
	var captcha map[string]string
	json.Unmarshal(c.Ctx.Input.RequestBody, &captcha)

	var configD = base64Captcha.ConfigDigit{
		Height:     42,
		Width:      89,
		MaxSkew:    beego.AppConfig.DefaultFloat("maxskew", 0.7),
		DotCount:   beego.AppConfig.DefaultInt("dotcount", 80),
		CaptchaLen: beego.AppConfig.DefaultInt("captchalen", 4),
	}

	idKeyD, capD := base64Captcha.GenerateCaptcha(captcha["captcha_id"], configD)
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)

	captcha_data := make(map[string]string)
	captcha_data["captcha_id"] = idKeyD
	captcha_data["captcha_val"] = base64stringD

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, captcha_data)
}

// 验证码验证
func (c *LoginController) VerfiyCaptcha(idkey, verifyValue string) bool {
	verifyResult := base64Captcha.VerifyCaptcha(idkey, verifyValue)
	if verifyResult {
		return true
	} else {
		return false
	}
}
