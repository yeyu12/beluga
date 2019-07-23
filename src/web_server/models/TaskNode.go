package models

import (
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
)

type TaskNode struct {
	Id                  int64     `json:"id"`
	Ip                  string    `json:"ip"`
	CreateTime          time.Time `json:"create_time"`
	Remake              string    `json:"remake"`
	IsDelete            int       `json:"is_delete"` // 0删除，1正常
}

func NewTaskNode() *TaskNode {
	return &TaskNode{}
}

// 加入表前缀的表名
func (m *TaskNode) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "task_node"
}

// 获取节点列表
func (m *TaskNode) List(o orm.Ormer, page int, page_size int, search string) helpers.Page {
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

	task_node_page_data := helpers.Page{
		Total:     total,
		TotalPage: math.Ceil(float64(total) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return task_node_page_data
}

// ID获取数据
func (m *TaskNode) IdFind(o orm.Ormer, id int64) TaskNode {
	var data TaskNode
	o.QueryTable(m.TableNamePrefix()).Filter("id", id).One(&data)

	return data
}

// 添加
func (m *TaskNode) Save(o orm.Ormer) (int64, error) {
	created, id, err := o.ReadOrCreate(m, "ip")

	if !created {
		m.Edit(o, map[string]interface{}{
			"is_delete": 1,
		})
	}

	return id, err
}

// 添加
func (m *TaskNode) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// ip获取数据
func (m *TaskNode) IpFind(o orm.Ormer, ip string) TaskNode {
	var data TaskNode
	o.QueryTable(m.TableNamePrefix()).Filter("ip", ip).One(&data)

	return data
}

// 修改
func (m *TaskNode) Edit(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("ip", m.Ip).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 根据id来修改
func (m *TaskNode) EditId(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// 删除
func (m *TaskNode) Del(o orm.Ormer, id int64) (int64, error) {
	res, err := o.QueryTable(m.TableNamePrefix()).Filter("id", id).Filter("is_delete", 0).Delete()

	return res, err
}

// id批量获取节点
func (m *TaskNode)IdInData(o orm.Ormer, ids interface{}) []orm.Params {
	var lists []orm.Params
	o.QueryTable(m.TableNamePrefix()).Filter("id__in", ids).Values(&lists)

	return lists
}