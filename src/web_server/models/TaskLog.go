package models

import (
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
)

type TaskLog struct {
	Id           int64     `json:"id"`
	TaskId       int64     `json:"task_id"`
	TaskExecType int       `json:"task_exec_type"`
	CreateTime   time.Time `json:"create_time"`
	EndTime      time.Time `json:"end_time"`
	NodeIp       string    `json:"node_ip"`
	ConsumeTime  float32   `json:"consume_time"`
	Err          string    `json:"err"`
	TaskName     string    `json:"task_name"`
	Cmd          string    `json:"cmd"`
	Output       string    `json:"output"`
}

// 当前对象
func NewTaskLog() *TaskLog {
	return &TaskLog{}
}

// 加入表前缀的表名
func (m *TaskLog) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "task_log"
}

func (m *TaskLog) List(o orm.Ormer, page int, page_size int, task_id string) helpers.Page {
	page_size_str := strconv.Itoa(page_size)
	total, _ := o.QueryTable(m.TableNamePrefix()).Filter("task_id", task_id).Count()

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select * from `" + m.TableNamePrefix() + "` where `task_id`=" + task_id
	sql += " order by id desc limit " + offset + "," + page_size_str

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	task_log_data := helpers.Page{
		Total:     total,
		TotalPage: math.Ceil(float64(total) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return task_log_data
}
