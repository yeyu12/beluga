package models

import (
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"time"
)

type Configuration struct {
	Id            int       `json:"id"`
	Appid         string    `json:"appid"`
	NamespaceName string    `json:"namespace_name"`
	Key           string    `json:"key"`
	Val           string    `json:"val"`
	Remake        string    `json:"remake"`
	CreateTime    time.Time `json:"create_time"`
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

// 加入表前缀的表名
func (m *Configuration) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration"
}

// 添加配置纪录
func (m *Configuration) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 更新配置记录
func (m *Configuration) Edit(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// appid和namespace_name和key更新配置
func (m *Configuration) AppidNamesapceNameKeyToEdit(o orm.Ormer, appid, namespace_name, key string, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).Filter("key", key).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 查找记录是否存在
func (m *Configuration) KeyFind(o orm.Ormer, appid, namespace_name, key string) Configuration {
	var config Configuration

	o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).
		Filter("key", key).One(&config)

	return config
}

// 删除
func (m *Configuration) Del(o orm.Ormer, appid, namespace_name, key string) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).
		Filter("namespace_name", namespace_name).Filter("key", key).Delete()

	return res, err
}

// 获取所有已发布的配置
func (m *Configuration) List(o orm.Ormer, appid, namespace_name string) []orm.Params {

	sql := "select `key`,`val`,`remake` from `" + m.TableNamePrefix() + "` where `appid`='" + appid +"' and namespace_name='" + namespace_name + "'"
	sql += " order by id desc "

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	return lists
}