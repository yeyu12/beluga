package models

import (
	"time"
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"strconv"
	"math"
)

type ConfigurationProject struct {
	Id          int       `json:"id"`
	ProjectName string    `json:"project_name"`
	AccountId   int       `json:"account_id"`
	CreateTime  time.Time `json:"create_time"`
	Appid       string    `json:"appid"`
}

func NewConfigurationProjectModel() *ConfigurationProject {
	return &ConfigurationProject{}
}

// 加入表前缀的表名
func (m *ConfigurationProject) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_project"
}

// 获取项目列表
func (m *ConfigurationProject) List(o orm.Ormer, page int, page_size int, search string) helpers.Page {
	account_table := helpers.GetTablePrefix() + "account"
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix())

	if search != "" {
		total_obj = total_obj.Filter("project_name__icontains", search)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select a.*,`b`.`nickname` from `" + m.TableNamePrefix() + "` a left join `" + account_table + "` `b` on a.account_id=b.id"
	if search != "" {
		sql += " where `a`.`project_name` like '%" + search + "%' "
	}
	sql += " order by `a`.`id` desc limit " + offset + "," + page_size_str

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

// 获取项目列表
func (m *ConfigurationProject) UserProjectList(o orm.Ormer, page int, page_size int, account_id int) helpers.Page {
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix())

	if account_id != 0 {
		total_obj = total_obj.Filter("account_id", account_id)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select * from `" + m.TableNamePrefix() + "`"
	if account_id != 0 {
		sql += " where `account_id`=" + strconv.Itoa(account_id)
	}
	sql += " order by id desc limit " + offset + "," + page_size_str

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

// 添加项目
func (m *ConfigurationProject) AddProject(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 修改项目
func (m *ConfigurationProject) EditProject(o orm.Ormer, data map[string]interface{}) bool {
	obj := o.QueryTable(m.TableNamePrefix())

	_, err := obj.Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 删除项目
func (m *ConfigurationProject) DelProject(o orm.Ormer, id int64) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("id", id).Delete()

	return res, err
}

// 项目id获取appid
func (m *ConfigurationProject) GetProjectIdToAppid(o orm.Ormer, project_id int) string {
	var config_project ConfigurationProject
	o.QueryTable(m.TableNamePrefix()).Filter("id", project_id).One(&config_project, "appid")

	return config_project.Appid
}

// 项目id获取项目数据
func (m *ConfigurationProject) GetProjectIdToData(o orm.Ormer, project_id int) ConfigurationProject {
	var config_project ConfigurationProject
	o.QueryTable(m.TableNamePrefix()).Filter("id", project_id).One(&config_project)

	return config_project
}

// 项目名称获取数据
func (m *ConfigurationProject) GetProjectNameToData(o orm.Ormer, name string) ConfigurationProject {
	var config_project ConfigurationProject
	o.QueryTable(m.TableNamePrefix()).Filter("project_name", name).One(&config_project)

	return config_project
}