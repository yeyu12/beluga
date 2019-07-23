package models

import (
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
)

type ConfigurationNode struct {
	Id                  int64     `json:"id"`
	Ip                  string    `json:"ip"`
	NodeConfId          string    `json:"node_conf_id"`
	ConfUpdateTime      time.Time `json:"conf_update_time"`
	CreateTime          time.Time `json:"create_time"`
	ConfUpdateAccountId int64     `json:"conf_update_account_id"`
	Remake              string    `json:"remake"`
	IsDelete            int       `json:"is_delete"` // 0删除，1正常
}

func NewConfigurationNode() *ConfigurationNode {
	return &ConfigurationNode{}
}

// 加入表前缀的表名
func (m *ConfigurationNode) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "configuration_node"
}

// 获取节点列表
func (m *ConfigurationNode) List(o orm.Ormer, page int, page_size int, search string) helpers.Page {
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix())

	if search != "" {
		total_obj = total_obj.Filter("ip__icontains", search)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select * from `" + m.TableNamePrefix() + "`"
	if search != "" {
		sql += " where `ip` like '%" + search + "%' "
	}
	sql += " order by id desc limit " + offset + "," + page_size_str

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	config_node_page_data := helpers.Page{
		Total:     total,
		TotalPage: math.Ceil(float64(total) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return config_node_page_data
}

// ID获取数据
func (m *ConfigurationNode) IdFind(o orm.Ormer, id int64) ConfigurationNode {
	var data ConfigurationNode
	o.QueryTable(m.TableNamePrefix()).Filter("id", id).One(&data)

	return data
}

// 添加
func (m *ConfigurationNode) Save(o orm.Ormer) (int64, error) {
	created, id, err := o.ReadOrCreate(m, "ip")

	if !created {
		m.Edit(o, map[string]interface{}{
			"is_delete": 1,
		})
	}

	return id, err
}

// 添加
func (m *ConfigurationNode) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// ip获取数据
func (m *ConfigurationNode) IpFind(o orm.Ormer, ip string) ConfigurationNode {
	var data ConfigurationNode
	o.QueryTable(m.TableNamePrefix()).Filter("ip", ip).One(&data)

	return data
}

// 修改
func (m *ConfigurationNode) Edit(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("ip", m.Ip).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 根据id来修改
func (m *ConfigurationNode) EditId(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 删除
func (m *ConfigurationNode) Del(o orm.Ormer, id int64) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("id", id).Filter("is_delete", 0).Delete()

	return res, err
}

// 使用该节点配置的节点ip
func (m *ConfigurationNode) NodeConfIdToIp(o orm.Ormer, node_conf_id string) []orm.Params {
	var lists []orm.Params
	sql := "select `ip` from " + m.TableNamePrefix() + " where is_delete=1 and node_conf_id like '%" + node_conf_id +"%'"
	o.Raw(sql).Values(&lists)

	return lists
}

// 查询在线节点数
func (m *ConfigurationNode) NodeConfIdOnLineCount(o orm.Ormer, node_conf_id string) []orm.Params {
	var lists []orm.Params
	sql := "select count(`ip`) as count_ip from " + m.TableNamePrefix() + " where is_delete=1 and node_conf_id like '%" + node_conf_id +"%'"
	o.Raw(sql).Values(&lists)

	return lists
}