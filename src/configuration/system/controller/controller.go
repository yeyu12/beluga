package controller

import (
	"net/http"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/jinzhu/gorm"
	goredis "github.com/go-redis/redis"
	"beluga/src/configuration/system/drive"
	"beluga/src/beluga/library"
)

type Controller struct {
	Conf    *drive.Config   // 配置
	Redis   *goredis.Client // redis链接
	Monitor *drive.Monitor  // 监控数据
	Writer  http.ResponseWriter
	Request *http.Request
	Mysql   *gorm.DB // mysql链接
	Etcd    library.EtcdClient
}

// http请求
func (c *Controller) Requests(writer http.ResponseWriter, request *http.Request) {
	c.Writer = writer
	c.Request = request
}

// 获取get参数
func (c *Controller) GetData() map[string]interface{} {
	get_data := make(map[string]interface{})
	request_get_data := c.Request.URL.Query()

	for key, val := range request_get_data {
		get_data[key] = val[0]
	}

	return get_data
}

// 获取post参数
func (c *Controller) PostData() map[string]interface{} {
	post_data := make(map[string]interface{})
	c.Request.ParseMultipartForm(32 << 20)
	if c.Request.MultipartForm != nil {
		for key, val := range c.Request.MultipartForm.Value {
			post_data[key] = val[0]
		}
	}

	return post_data
}

// 返回json数据
func (c *Controller) Json(code int, errOk string, msg string, data ...interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(code)

	res := map[string]interface{}{
		"err":  errOk,
		"msg":  msg,
		"data": data[0],
	}

	res_json, err := json.Marshal(res)
	if err != nil {
		drive.Notices(res, errors.Wrap(err, "数据返回失败"))
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)

		return
	}

	c.Writer.Write(res_json)
	c.StopRun()
}

// 中断流程继续
func (c *Controller) StopRun() {
	panic(errors.New("User err"))
}
