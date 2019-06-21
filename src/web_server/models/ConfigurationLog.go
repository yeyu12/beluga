package models

import (
	"time"
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
)

type ConfigurationLog struct {
	Id            int       `json:"id"`
	Appid         string    `json:"appid"`
	NamespaceName string    `json:"namespace_name"`
	Version       string    `json:"version"`
	AccountId     int       `json:"account_id"`
	Key           string    `json:"key"`
	Val           string    `json:"val"`
	Remake        string    `json:"remake"`
	Type          int       `json:"type"`
	IsRelease     int       `json:"is_release"`
	OldVal        string    `json:"old_val"`
	OldRemake     string    `json:"old_remake"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}

func NewConfigurationLog() *ConfigurationLog {
	return &ConfigurationLog{}
}

// 加入表前缀的表名
func (m *ConfigurationLog) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_log"
}

// 添加配置纪录
func (m *ConfigurationLog) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 更新配置记录
func (m *ConfigurationLog) Edit(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// appid、namespace_name、version更新数据
func (m *ConfigurationLog) AppidNamespaceNameVersionToEdit(o orm.Ormer, appid, namespace_name string, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 查找记录是否存在
func (m *ConfigurationLog) KeyFind(o orm.Ormer, appid, namespace_name, version, key string, release int) ConfigurationLog {
	var config_log ConfigurationLog

	obj := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).
		Filter("key", key).Filter("is_release", release)
	if version != "" {
		obj = obj.Filter("version", version)
	}
	obj.One(&config_log)

	return config_log
}

// 删除
func (m *ConfigurationLog) Del(o orm.Ormer, appid, namespace_name, key string) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).
		Filter("namespace_name", namespace_name).Filter("key", key).Filter("is_release", 0).Delete()

	return res, err
}

// 删除版本号下搜索未发布的配置
func (m *ConfigurationLog)DelLog(o orm.Ormer, appid, namespace_name, version string) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).
		Filter("namespace_name", namespace_name).Filter("version", version).Filter("is_release", 0).Delete()

	return res, err
}

// 版本号获取未发布的配置日志
func (m *ConfigurationLog) VersionToConfigLog(o orm.Ormer, appid, namespace_name, version string) []orm.Params {
	sql := "select * from `" + m.TableNamePrefix() + "` where `appid`='" + appid + "' and `namespace_name`='" + namespace_name + "' and `version`='" + version + "'"
	var config_log []orm.Params
	o.Raw(sql).Values(&config_log)

	return config_log
}

// 版本号获取发布的配置日志
func (m *ConfigurationLog) VersionToConfigToReleaseLog(o orm.Ormer, appid, namespace_name, version string) []orm.Params {
	sql := "select * from `" + m.TableNamePrefix() + "` where `appid`='" + appid + "' and `namespace_name`='" + namespace_name + "' and `version`='" + version + "'"
	sql += " and `is_release`=1 order by id desc"
	var config_log []orm.Params
	o.Raw(sql).Values(&config_log)

	return config_log
}

// 获取所有未发布的配置
func (m *ConfigurationLog) List(o orm.Ormer, appid, namespace_name string) []orm.Params {
	account_table := helpers.GetTablePrefix() + "account"
	sql := "select a.`key`,`a`.`val`,`a`.`remake`,`a`.`type`,`a`.`update_time`,`b`.`username` from `" +
		m.TableNamePrefix() + "` a left join `" + account_table + "` b on `b`.`id`=`a`.`account_id` where `a`.`appid`='" + appid + "' and `a`.`namespace_name`='" + namespace_name + "'"
	sql += " and `a`.`is_release`=0 order by `a`.`id` desc"

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	return lists
}

// appid和namespace_name和key更新配置
func (m *ConfigurationLog) AppidNamesapceNameKeyToEdit(o orm.Ormer, appid, namespace_name, key string, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).Filter("key", key).Update(data)

	if err == nil {
		return true
	}

	return false
}