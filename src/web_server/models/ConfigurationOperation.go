package models

import (
	"beluga/src/web_server/helpers"
	"time"
	"github.com/astaxie/beego/orm"
)

type ConfigurationOperation struct {
	Id            int       `json:"id"`
	Appid         string    `json:"appid"`
	NamespaceName string    `json:"namespace_name"`
	Version       string    `json:"version"`
	AccountId     int       `json:"account_id"`
	OperationType int       `json:"operation_type"`
	Rollback      int       `json:"rollback"`
	CreateTime    time.Time `json:"create_time"`
}

func NewConfigurationOperation() *ConfigurationOperation {
	return &ConfigurationOperation{}
}

// 加入表前缀的表名
func (m *ConfigurationOperation) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_operation"
}

// 获取操作纪录
func (m *ConfigurationOperation) List(o orm.Ormer, appid, namespace_name string) []orm.Params {
	account_table := helpers.GetTablePrefix() + "account"
	sql := "select a.*,`b`.`username` as `account_name` from `" + m.TableNamePrefix() + "` `a` " +
		"left join `" + account_table + "` b on `a`.`account_id`=`b`.`id` where `appid`='" + appid + "' " +
		"and `namespace_name`='" + namespace_name + "'"

	sql += " order by id desc"
	var list []orm.Params
	o.Raw(sql).Values(&list)

	return list
}

// 添加操作记录
func (m *ConfigurationOperation) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 获取上次发布的数据
func (m *ConfigurationOperation) LastRecords(o orm.Ormer, appid, namespace_name string) ConfigurationOperation {
	var data ConfigurationOperation
	o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).
		Filter("operation_type", 0).Filter("rollback", 0).OrderBy("-version").One(&data)

	return data
}

// 修改操作是为回滚
func (m *ConfigurationOperation) EditOperationRollback(o orm.Ormer, appid, namespace_name, version string) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("appid", appid).Filter("namespace_name", namespace_name).
		Filter("version", version).Update(map[string]interface{}{
		"rollback": -1,
	})

	if err == nil {
		return true
	}

	return false
}
