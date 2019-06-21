package models

import (
	"time"
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
)

type ConfigurationVersion struct {
	Id            int       `json:"id"`
	PId           int       `json:"p_id"`
	Version       string    `json:"version"`
	NamespaceName string    `json:"namespace_name"`
	Appid         string    `json:"appid"`
	IsRelease     int       `json:"is_release"`
	CreateTime    time.Time `json:"create_time"`
}

func NewConfigurationVersion() *ConfigurationVersion {
	return &ConfigurationVersion{}
}

// 加入表前缀的表名
func (m *ConfigurationVersion) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_version"
}

// 获取最新未发布的版本号
func (m *ConfigurationVersion) GetVersion(o orm.Ormer, appid, namespace_name string) ConfigurationVersion {
	var config_version ConfigurationVersion

	o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).
		Filter("namespace_name", namespace_name).Filter("is_release", 0).
		OrderBy("-id").One(&config_version)

	return config_version
}

// 获取已发布的版本
func (m *ConfigurationVersion) GetReleaseVersion(o orm.Ormer, appid, namespace_name string) ConfigurationVersion {
	var config_release_version ConfigurationVersion
	o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).
		Filter("namespace_name", namespace_name).Filter("is_release", 1).
		OrderBy("-id").One(&config_release_version)

	return config_release_version
}

// 添加版本
func (m *ConfigurationVersion) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 更新
func (m *ConfigurationVersion) Update(o orm.Ormer, data map[string]interface{}) bool {
	obj := o.QueryTable(m.TableNamePrefix())

	_, err := obj.Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}