package models

import (
	"beluga/src/web_server/helpers"
	"time"
	"github.com/astaxie/beego/orm"
)

type ConfigurationOperationLog struct {
	Id                       int       `json:"id"`
	ConfigurationOperationId int       `json:"configuration_operation_id"`
	Key                      string    `json:"key"`
	Val                      string    `json:"val"`
	AccountId                int       `json:"account_id"`
	Remake                   string    `json:"remake"`
	CreateTime               time.Time `json:"create_time"`
}

func NewConfigurationOperationLog() *ConfigurationOperationLog {
	return &ConfigurationOperationLog{}
}

// 加入表前缀的表名
func (m *ConfigurationOperationLog) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_operation_log"
}

// 添加操作记录日志
func (m *ConfigurationOperationLog) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 版本号获取发布的配置日志
func (m *ConfigurationOperationLog) VersionToConfigToReleaseLog(o orm.Ormer, appid, namespace_name, version string) []orm.Params {
	config_operation_table := helpers.GetTablePrefix() + "configuration_operation"

	sql := "select a.* from `" + m.TableNamePrefix() + "` `a` right join `" +
		config_operation_table + "` `b` on `b`.`appid`='" + appid + "' and `b`.`namespace_name`='" + namespace_name + "'" +
		" and `b`.`version`='" + version + "' where `b`.`id`=`a`.`configuration_operation_id`"

	var config_log []orm.Params
	o.Raw(sql).Values(&config_log)

	return config_log
}
