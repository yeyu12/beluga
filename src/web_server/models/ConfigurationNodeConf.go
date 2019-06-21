package models

import (
	"time"
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"strconv"
	"math"
)

type ConfigurationNodeConf struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Conf       string    `json:"conf"`
	AccountId  int       `json:"account_id"`
	CreateTime time.Time `json:"create_time"`
}

func NewConfigurationNodeConf() *ConfigurationNodeConf {
	return &ConfigurationNodeConf{}
}

// 加入表前缀的表名
func (m *ConfigurationNodeConf) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_node_conf"
}

// 配置名称查找数据
func (m *ConfigurationNodeConf) NameFind(o orm.Ormer, name string) []orm.Params {
	account_table := helpers.GetTablePrefix() + "account"
	sql := "select a.*,`b`.`username` from `" + m.TableNamePrefix() + "` `a` " +
		"left join `" + account_table + "` b on `a`.`account_id`=`b`.`id` where `name`='" + name + "'"
	sql += " order by id desc"

	var list []orm.Params
	o.Raw(sql).Values(&list)

	return list
}

// 添加
func (m *ConfigurationNodeConf) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 修改
func (m *ConfigurationNodeConf) Edit(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 删除
func (m *ConfigurationNodeConf) Del(o orm.Ormer, id int64) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("id", id).Delete()

	return res, err
}

// 获取项目列表
func (m *ConfigurationNodeConf) List(o orm.Ormer, page int, page_size int, search string) helpers.Page {
	account_table := helpers.GetTablePrefix() + "account"
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix())

	if search != "" {
		total_obj = total_obj.Filter("name__icontains", search)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select a.*,b.nickname from `" + m.TableNamePrefix() + "` `a` " +
		" left join `" + account_table + "` `b` on `a`.`account_id`=`b`.`id`"

	if search != "" {
		sql += " where `a`.`name` like '%" + search + "%' "
	}
	sql += " order by `a`.id desc limit " + offset + "," + page_size_str

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	config_project_page_data := helpers.Page{
		Total:     total,
		TotalPage: math.Ceil(float64(total) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return config_project_page_data
}

// 批量ID查询数据
func (m *ConfigurationNodeConf) IdsToData(o orm.Ormer, ids string) []orm.Params {
	var lists []orm.Params
	sql := "select * from " + m.TableNamePrefix() + " where id in (" + ids +")"
	o.Raw(sql).Values(&lists)

	return lists
}