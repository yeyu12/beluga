package models

import (
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
)

type OperationLog struct {
	Id         int64     `json:"id"`
	C          string    `json:"c"`
	Params     string    `json:"params"`
	AccountId  int64     `json:"account_id"`
	Ident      string    `json:"ident"`
	CreateTime time.Time `json:"create_time"`
}

func NewOperationLog() *OperationLog {
	return &OperationLog{}
}

// 加入表前缀的表名
func (m *OperationLog) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "operation_log"
}

// 获取操作列表
func (m *OperationLog) List(o orm.Ormer, page int, page_size int, account_id int64) helpers.Page {
	account_table := helpers.GetTablePrefix() + "account"
	page_size_str := strconv.Itoa(page_size)

	total_obj := o.QueryTable(m.TableNamePrefix())

	if account_id != 0 {
		total_obj = total_obj.Filter("account_id", account_id)
	}

	total, _ := total_obj.Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select a.*,b.nickname from `" + m.TableNamePrefix() + "` a"
	sql += " left join " + account_table + " b on a.account_id=b.id"
	if account_id != 0 {
		account_id_str := strconv.Itoa(int(account_id))
		sql += " where a.`account_id` =" + account_id_str + " "
	}
	sql += " order by a.id desc limit " + offset + "," + page_size_str

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

// 添加操作记录日志
func (m *OperationLog) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}
