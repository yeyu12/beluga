package models

import (
	"github.com/astaxie/beego"
	"time"
	"strconv"
	"github.com/astaxie/beego/orm"
	"math"
	"beluga/src/web_server/helpers"
)

type ConfigurationNamespace struct {
	Id            int       `json:"id"`
	ProjectId     int       `json:"project_id"`
	NamespaceName string    `json:"namespace_name"`
	AccountId     int       `json:"account_id"`
	CreateTime    time.Time `json:"create_time"`
}

func NewConfigurationNamespaceModel() *ConfigurationNamespace {
	return &ConfigurationNamespace{}
}

// 加入表前缀的表名
func (m *ConfigurationNamespace) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_namespace"
}

func (m *ConfigurationNamespace) List(o orm.Ormer, page int, page_size int, project_id int, account_id int, search string) helpers.Page {
	account_table := helpers.GetTablePrefix() + "account"
	config_project_table := helpers.GetTablePrefix() + "configuration_project"
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix()).Filter("project_id", project_id)

	if search != "" {
		total_obj = total_obj.Filter("namespace_name__icontains", search)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select a.*,b.nickname,c.project_name from `" + m.TableNamePrefix() + "` a left join `" + account_table + "` b on b.id=a.account_id"
	sql += " left join `" + config_project_table + "` c on c.id=a.project_id"
	sql += " where `a`.`project_id`=" + strconv.Itoa(project_id)
	if search != "" {
		sql += " and `a`.`namespace_name` like '%" + search + "%' "
	}
	sql += " order by `a`.`id` desc limit " + offset + "," + page_size_str
	beego.Error(sql)

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

// 添加命名空间
func (m *ConfigurationNamespace) AddNamespace(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 修改命名空间
func (m *ConfigurationNamespace) EditNamespace(o orm.Ormer, data map[string]interface{}) bool {
	obj := o.QueryTable(m.TableNamePrefix())

	_, err := obj.Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 删除命名空间
func (m *ConfigurationNamespace) DelNamespace(o orm.Ormer, id int64) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("id", id).Delete()

	return res, err
}

// 命名空间id获取命名空间名称
func (m *ConfigurationNamespace) GetNamespaceIdToName(o orm.Ormer, project_id, namespace_id int) string {
	var config_namespace ConfigurationNamespace
	o.QueryTable(m.TableNamePrefix()).Filter("id", namespace_id).Filter("project_id", project_id).One(&config_namespace, "namespace_name")

	return config_namespace.NamespaceName
}

// 命名空间id获取命名空间数据
func (m *ConfigurationNamespace) GetNamespaceIdToData(o orm.Ormer, project_id, namespace_id int) ConfigurationNamespace {
	var config_namespace ConfigurationNamespace
	o.QueryTable(m.TableNamePrefix()).Filter("id", namespace_id).Filter("project_id", project_id).One(&config_namespace)

	return config_namespace
}

// 项目id和命名空间获取数据
func (m *ConfigurationNamespace) GetProjectIdNamespaceNameToData(o orm.Ormer, project_id int, namespace_name string) ConfigurationNamespace {
	var config_namespace ConfigurationNamespace
	o.QueryTable(m.TableNamePrefix()).Filter("namespace_name", namespace_name).Filter("project_id", project_id).One(&config_namespace)

	return config_namespace
}