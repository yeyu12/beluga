package models

import (
	"beluga/src/web_server/helpers"
	"github.com/astaxie/beego/orm"
	"math"
	"strconv"
	"time"
)

type Task struct {
	Id                int64     `json:"id"`
	Name              string    `json:"name"`
	CreateTime        time.Time `json:"create_time"`
	StartTime         time.Time `json:"start_time"`
	ConsumeTime       float64   `json:"consume_time"`
	Overtime          int       `json:"overtime"`
	LastExecType      int       `json:"last_exec_type"`
	TaskType          int       `json:"task_type"`
	Rely              int       `json:"rely"`
	SubtasksId        string    `json:"subtasks_id"`
	Cron              string    `json:"cron"`
	TaskExecType      string    `json:"task_exec_type"`
	ExecTaskNodeType  int       `json:"exec_task_node_type"`
	ExecTaskNodeId    string    `json:"exec_task_node_id"`
	Cmd               string    `json:"cmd"`
	HttpType          string    `json:"http_type"`
	TaskFailNum       int       `json:"task_fail_num"`
	TaskFailRetryTime int       `json:"task_fail_retry_time"`
	TaskNotice        int       `json:"task_notice"`
	NoticeType        int       `json:"notice_type"`
	KeywordNotice     string    `json:"keyword_notice"`
	Remake            string    `json:"remake"`
	AccountId         int       `json:"account_id"`
	Status            int       `json:"status"`
	NextExecTime      time.Time `json:"next_exec_time"`
}

// 当前对象
func NewTask() *Task {
	return &Task{}
}

// 加入表前缀的表名
func (m *Task) TableNamePrefix() string {
	return helpers.GetTablePrefix() + "task"
}

// 列表
func (m *Task) List(o orm.Ormer, page int, page_size int, search string, task_id string) helpers.Page {
	account_table := helpers.GetTablePrefix() + "account"
	page_size_str := strconv.Itoa(page_size)

	total_sql := "SELECT COUNT(*) as `count` FROM `" + m.TableNamePrefix() + "` WHERE `status` !=-1"

	if search != "" {
		total_sql += " and name like '%" + search + "%' "
	}
	if task_id != "" && task_id != "0" {
		total_sql += " and id=" + task_id
	}

	var total helpers.Total
	_ = o.Raw(total_sql).QueryRow(&total)

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select a.*,b.nickname from `" + m.TableNamePrefix() + "` `a`"
	sql += " left join `" + account_table + "` `b` on b.id=a.account_id"
	sql += " where `a`.`status`!=-1"

	if search != "" {
		sql += " and `a`.`name` like '%" + search + "%' "
	}
	if task_id != "" && task_id != "0" {
		sql += " and `a`.`id`=" + task_id
	}

	sql += " order by `a`.`id` desc limit " + offset + "," + page_size_str

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	task_page_data := helpers.Page{
		Total:     total.Count,
		TotalPage: math.Ceil(float64(total.Count) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return task_page_data
}

// 子任务数据
func (m *Task) SubtasksList(o orm.Ormer, page int, page_size int, search string) helpers.Page {
	account_table := helpers.GetTablePrefix() + "account"
	page_size_str := strconv.Itoa(page_size)

	total_sql := "SELECT COUNT(*) as `count` FROM `" + m.TableNamePrefix() + "` WHERE `status` !=-1 AND `task_type` = 2"

	if search != "" {
		total_sql += " and name like '%" + search + "%' "
	}

	var total helpers.Total
	_ = o.Raw(total_sql).QueryRow(&total)

	offset := strconv.Itoa(((page - 1) * page_size))
	sql := "select a.*,b.nickname from `" + m.TableNamePrefix() + "` `a`"
	sql += " left join `" + account_table + "` `b` on b.id=a.account_id"
	sql += " where `a`.`status`!=-1 and `a`.`task_type`=2"

	if search != "" {
		sql += " and `a`.`name` like '%" + search + "%' "
	}

	sql += " order by `a`.`id` desc limit " + offset + "," + page_size_str

	var lists []orm.Params
	o.Raw(sql).Values(&lists)

	task_page_data := helpers.Page{
		Total:     total.Count,
		TotalPage: math.Ceil(float64(total.Count) / float64(page_size)),
		PageSize:  page_size,
		Page:      page,
		List:      lists,
	}

	return task_page_data
}

// id获取数据
func (m *Task) TaskIdToData(o orm.Ormer, task_id int64) (*Task, error) {
	task := &Task{}
	err := o.QueryTable(m.TableNamePrefix()).Filter("id", task_id).One(task)

	return task, err
}

// 添加数据
func (m *Task) Add(o orm.Ormer) (int64, error) {
	id, err := o.Insert(m)

	return id, err
}

// 更新配置记录
func (m *Task) Edit(o orm.Ormer, data map[string]interface{}) bool {
	_, err := o.QueryTable(m.TableNamePrefix()).Filter("id", m.Id).Update(data)

	if err == nil {
		return true
	}

	return false
}

// id批量获取任务
func (m *Task)IdInData(o orm.Ormer, id interface{}) []orm.Params {
	var lists []orm.Params
	o.QueryTable(m.TableNamePrefix()).Filter("id__in", id).Values(&lists)

	return lists
}