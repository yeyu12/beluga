package controllers

import (
	"encoding/json"
	"io"
	"time"
	"strconv"
	"strings"
	"runtime"
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BaseController struct {
	beego.Controller
	AccountInfo *models.Account
	PageSize    int
	Orm         orm.Ormer
}

// 初始化
func (c *BaseController) Prepare() {
	c.AllowCross()

	c.PageSize = beego.AppConfig.DefaultInt("page_size", 10)

	token := c.Ctx.Request.Header.Get("Account-token")

	auth_arr := []string{
		"/login",
	}

	c.Orm = orm.NewOrm()

	if token == "" {
		if !helpers.In_array(c.Ctx.Request.RequestURI, auth_arr) && (c.Ctx.Request.RequestURI != "/") {
			c.ResponseJson(helpers.LOGIN_FATL_CODE, false, helpers.LOGIN_FATL_MSG)
		}
	} else {
		c.TokenToUserInfo(token)
		c.TokenIsExpiry()

		if c.AccountInfo.Status != 1 {
			c.ResponseJson(helpers.USER_PROHIBIT_CODE, false, helpers.USER_PROHIBIT_MSG)
		}
	}
}

// AllowCross 跨域
func (c *BaseController) AllowCross() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                                                          //允许访问源
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")                                   //允许post访问
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With,Account-token") //header的类型
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("Content-type", "application/json;charset=utf-8") //返回数据格式是json TODO 在开发的时候显示json
	c.Ctx.ResponseWriter.Header().Set("Cache-Control", "no-cache, no-store")
}

// ResponseJson json返回
func (c *BaseController) ResponseJson(code int, status bool, msg string, data ...interface{}) {
	resMap := make(map[string]interface{}, 3)

	resMap["code"] = code
	resMap["msg"] = msg
	resMap["status"] = status

	data_len := len(data)

	if (data_len > 0) && data[0] != nil {
		if data_len > 1 {
			resMap["data"] = data
		} else {
			resMap["data"] = data[0]
		}
	}

	resJson, err := json.Marshal(resMap)

	if err != nil {
		beego.Error(err)
	}

	c.recordOperationLog()

	c.AllowCross()
	io.WriteString(c.Ctx.ResponseWriter, string(resJson))
	c.StopRun()
}

// 返回string
func (c *BaseController) ResponseStr(str string) {
	c.AllowCross()
	io.WriteString(c.Ctx.ResponseWriter, str)
	c.StopRun()
}

func (c *BaseController) ResStr(code int, status bool, msg string, data ...interface{}) string {
	resMap := make(map[string]interface{}, 3)

	resMap["code"] = code
	resMap["msg"] = msg
	resMap["status"] = status

	data_len := len(data)

	if (data_len > 0) && data[0] != nil {
		if data_len > 1 {
			resMap["data"] = data
		} else {
			resMap["data"] = data[0]
		}
	}

	resJson, err := json.Marshal(resMap)

	if err != nil {
		beego.Error(err)
	}

	return string(resJson)
}

// token 获取用户数据
func (c *BaseController) TokenToUserInfo(token string) {
	c.AccountInfo = models.NewAccount()
	account, err := c.AccountInfo.TokenToUser(c.Orm, token)

	if err != nil {
		c.ResponseJson(helpers.TOKEN_FATL_CODE, false, helpers.TOKEN_FATL_MSG)
	}

	c.AccountInfo = account
}

// token 是否过期
func (c *BaseController) TokenIsExpiry() {
	current_time := time.Now().Unix()
	token_effective := c.AccountInfo.TokenExpiryTime.Unix()
	token_expired_time := beego.AppConfig.DefaultInt64("token_expired_time", 1296000)

	if current_time > (token_effective + token_expired_time) {
		c.ResponseJson(helpers.TOKEN_FATL_CODE, false, helpers.TOKEN_FATL_MSG)
	}
}

// 获取系统信息
func (c *BaseController) GetSystemInfo() {
	sys_data := map[string]string{
		"year":             strconv.Itoa(time.Now().Year()),
		"golang_version":   strings.ToUpper(runtime.Version()),
		"appname":          beego.AppConfig.DefaultString("appname", "beluga"),
		"system_version":   beego.AppConfig.DefaultString("system_version", "0.0.1"),
		"official_network": beego.AppConfig.DefaultString("official_network", "http://www.beluga.org.cn"),
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, sys_data)
}

func (c *BaseController) Options() {
	c.AllowCross()
	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 记录操作
func (c *BaseController) recordOperationLog()  {
	if c.Ctx.Request.Method != "OPTIONS"{
		operation_log_model := models.NewOperationLog()
		operation_log_model.C = c.Ctx.Request.RequestURI
		operation_log_model.Ident = c.Ctx.Request.Header.Get("User-Agent")
		operation_log_model.Params = string(c.Ctx.Input.RequestBody)
		operation_log_model.CreateTime = time.Now()

		if c.AccountInfo != nil {
			operation_log_model.AccountId = int64(c.AccountInfo.Id)
		}

		operation_log_model.Add(c.Orm)
	}
}
