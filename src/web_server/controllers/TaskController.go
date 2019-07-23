package controllers

import (
	"beluga/src/beluga/library"
	"beluga/src/beluga/task_constant"
	"beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorhill/cronexpr"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TaskController struct {
	BaseController
}

// 列表
func (c *TaskController) List() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	page, _ := strconv.Atoi(request_data["page"])
	search := request_data["search"]
	task_id := request_data["task_id"]

	if page == 0 {
		page = 1
	}

	task_model := models.NewTask()
	data := task_model.List(c.Orm, page, c.PageSize, search, task_id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 获取子任务列表
func (c *TaskController) SubtasksList() {
	request_data := make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	search := request_data["search"]
	page, _ := strconv.Atoi(request_data["page"])

	if page == 0 {
		page = 1
	}

	task_model := models.NewTask()
	data := task_model.SubtasksList(c.Orm, page, c.PageSize, search)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, data)
}

// 任务id获取任务数据
func (c *TaskController) GetTaskIdToInfo() {
	request_data := make(map[string]float64)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["task_id"] == 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	task_model := models.NewTask()
	task_data, err := task_model.TaskIdToData(c.Orm, int64(request_data["task_id"]))
	if err != nil {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, task_data)
}

// 添加
func (c *TaskController) Add() {
	request_data := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	// 数据验证
	cron_obj, err := cronexpr.Parse(request_data["cron"].(string));
	if err != nil {
		c.ResponseJson(helpers.TASK_CRONTAB_ANALYSIS_FAIL_CODE, false, helpers.TASK_CRONTAB_ANALYSIS_FAIL_MSG)
	}
	if request_data["name"].(string) == "" {
		c.ResponseJson(helpers.TASK_NAME_FAIL_CODE, false, helpers.TASK_NAME_FAIL_MSG)
	}
	if request_data["task_type"].(string) == "" {
		request_data["task_type"] = "1"
	}
	if request_data["task_type"] == "1" {
		if request_data["rely"].(string) == "" {
			request_data["rely"] = "1"
		}

		subtasks_id := request_data["subtasks_id"]
		delete(request_data, "subtasks_id")
		request_data["subtasks_id"] = strings.Replace(strings.Trim(fmt.Sprint(subtasks_id), "[]"), " ", ",", -1)

		if request_data["subtasks_id"] == "" {
			request_data["subtasks_id"] = ""
		}
	}
	if request_data["task_exec_type"].(string) == "" {
		request_data["task_exec_type"] = "shell"
	}
	if request_data["task_exec_type"] == "shell" {
		if request_data["http_type"].(string) == "" {
			request_data["http_type"] = "post"
		}
	}

	exec_task_node_id := request_data["exec_task_node_id"]
	delete(request_data, "exec_task_node_id")
	request_data["exec_task_node_id"] = strings.Replace(strings.Trim(fmt.Sprint(exec_task_node_id), "[]"), " ", ",", -1)
	if request_data["exec_task_node_id"] == "" {
		request_data["exec_task_node_id"] = "0"
	}

	if request_data["cmd"].(string) == "" {
		c.ResponseJson(helpers.TASK_CMD_FAIL_CODE, false, helpers.TASK_CMD_FAIL_MSG)
	}
	overtime_int := int(request_data["overtime"].(float64))
	if overtime_int < 0 || overtime_int > 86400 {
		c.ResponseJson(helpers.TASK_OVERTIME_FAIL_CODE, false, helpers.TASK_OVERTIME_FAIL_MSG)
	}
	task_fail_num_int := int(request_data["task_fail_num"].(float64))
	if task_fail_num_int < 0 {
		c.ResponseJson(helpers.TASK_FAIL_NUM_CODE, false, helpers.TASK_FAIL_NUM_MSG)
	}
	task_fail_retry_time_int := int(request_data["task_fail_retry_time"].(float64))
	if task_fail_retry_time_int < 0 {
		c.ResponseJson(helpers.TASK_FAIL_RETRY_TIME_CODE, false, helpers.TASK_FAIL_RETRY_TIME_MSG)
	}
	if request_data["task_notice"].(string) == "" {
		request_data["task_notice"] = "0"
	} else {
		task_notice_int, _ := strconv.Atoi(request_data["task_notice"].(string))
		if task_notice_int < 0 || task_notice_int > 3 {
			c.ResponseJson(helpers.TASK_NOTICE_FAIL_CODE, false, helpers.TASK_NOTICE_FAIL_MSG)
		}
	}
	if request_data["task_notice"] == "3" {
		if request_data["keyword_notice"].(string) == "" {
			c.ResponseJson(helpers.TASK_KEYWORD_NOTICE_FAIL_CODE, false, helpers.TASK_KEYWORD_NOTICE_FAIL_MSG)
		}
	}

	beego.Error(reflect.TypeOf(request_data["subtasks_id"]))
	// 保存数据
	task_model := models.NewTask()
	task_model.Name = request_data["name"].(string)
	task_model.CreateTime = time.Now()
	task_model.Overtime = int(request_data["overtime"].(float64))
	task_model.TaskType, _ = strconv.Atoi(request_data["task_type"].(string))
	task_model.Rely, _ = strconv.Atoi(request_data["rely"].(string))
	task_model.Status = 1
	if request_data["task_type"] == "1" {
		task_model.SubtasksId = request_data["subtasks_id"].(string)
	} else {
		task_model.Status = 2
	}
	task_model.Cron = request_data["cron"].(string)
	task_model.TaskExecType = request_data["task_exec_type"].(string)
	task_model.ExecTaskNodeId = request_data["exec_task_node_id"].(string)
	task_model.Cmd = request_data["cmd"].(string)
	task_model.HttpType = request_data["http_type"].(string)
	task_model.TaskFailNum = int(request_data["task_fail_num"].(float64))
	task_model.TaskFailRetryTime = int(request_data["task_fail_retry_time"].(float64))
	task_model.TaskNotice, _ = strconv.Atoi(request_data["task_notice"].(string))
	task_model.NoticeType, _ = strconv.Atoi(request_data["notice_type"].(string))
	task_model.KeywordNotice = request_data["keyword_notice"].(string)
	task_model.Remake = request_data["remake"].(string)
	task_model.AccountId = c.AccountInfo.Id
	task_model.NextExecTime = cron_obj.Next(time.Now())

	if task_model.ExecTaskNodeId == "0" {
		task_model.ExecTaskNodeType = 1
	} else {
		task_model.ExecTaskNodeType = 0
	}

	_, err = task_model.Add(c.Orm)
	if err != nil {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	if task_model.Status == 1 {
		c.addTaskEtcd(task_model, false)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, request_data)
}

// 修改
func (c *TaskController) Edit() {
	request_data := make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	// 数据验证
	cron_obj, err := cronexpr.Parse(request_data["cron"].(string));
	if err != nil {
		c.ResponseJson(helpers.TASK_CRONTAB_ANALYSIS_FAIL_CODE, false, helpers.TASK_CRONTAB_ANALYSIS_FAIL_MSG)
	}
	if request_data["name"].(string) == "" {
		c.ResponseJson(helpers.TASK_NAME_FAIL_CODE, false, helpers.TASK_NAME_FAIL_MSG)
	}
	if request_data["task_type"].(string) == "" {
		request_data["task_type"] = "1"
	}
	if request_data["task_type"] == "1" {
		if request_data["rely"].(string) == "" {
			request_data["rely"] = "1"
		}

		subtasks_id := request_data["subtasks_id"]
		delete(request_data, "subtasks_id")
		request_data["subtasks_id"] = strings.Replace(strings.Trim(fmt.Sprint(subtasks_id), "[]"), " ", ",", -1)

		if request_data["subtasks_id"] == "" {
			request_data["subtasks_id"] = "0"
		}
	}
	if request_data["task_exec_type"].(string) == "" {
		request_data["task_exec_type"] = "shell"
	}
	if request_data["task_exec_type"] == "shell" {
		if request_data["http_type"].(string) == "" {
			request_data["http_type"] = "post"
		}
	}

	exec_task_node_id := request_data["exec_task_node_id"]
	delete(request_data, "exec_task_node_id")
	request_data["exec_task_node_id"] = strings.Replace(strings.Trim(fmt.Sprint(exec_task_node_id), "[]"), " ", ",", -1)
	if request_data["exec_task_node_id"] == "" {
		request_data["exec_task_node_id"] = "0"
	}

	if request_data["cmd"].(string) == "" {
		c.ResponseJson(helpers.TASK_CMD_FAIL_CODE, false, helpers.TASK_CMD_FAIL_MSG)
	}
	overtime_int := int(request_data["overtime"].(float64))
	if overtime_int < 0 || overtime_int > 86400 {
		c.ResponseJson(helpers.TASK_OVERTIME_FAIL_CODE, false, helpers.TASK_OVERTIME_FAIL_MSG)
	}
	task_fail_num_int := int(request_data["task_fail_num"].(float64))
	if task_fail_num_int < 0 {
		c.ResponseJson(helpers.TASK_FAIL_NUM_CODE, false, helpers.TASK_FAIL_NUM_MSG)
	}
	task_fail_retry_time_int := int(request_data["task_fail_retry_time"].(float64))
	if task_fail_retry_time_int < 0 {
		c.ResponseJson(helpers.TASK_FAIL_RETRY_TIME_CODE, false, helpers.TASK_FAIL_RETRY_TIME_MSG)
	}
	if request_data["task_notice"].(string) == "" {
		request_data["task_notice"] = "0"
	} else {
		task_notice_int, _ := strconv.Atoi(request_data["task_notice"].(string))
		if task_notice_int < 0 || task_notice_int > 3 {
			c.ResponseJson(helpers.TASK_NOTICE_FAIL_CODE, false, helpers.TASK_NOTICE_FAIL_MSG)
		}
	}
	if request_data["task_notice"] == "3" {
		if request_data["keyword_notice"].(string) == "" {
			c.ResponseJson(helpers.TASK_KEYWORD_NOTICE_FAIL_CODE, false, helpers.TASK_KEYWORD_NOTICE_FAIL_MSG)
		}
	}
	task_id := int64(request_data["id"].(float64))
	if task_id == 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 保存数据
	task_model := models.NewTask()
	task_update_data := make(map[string]interface{})
	task_model.Id = task_id
	task_update_data["name"] = request_data["name"].(string)
	task_update_data["overtime"] = int(request_data["overtime"].(float64))
	task_update_data["task_type"], _ = strconv.Atoi(request_data["task_type"].(string))
	task_update_data["rely"], _ = strconv.Atoi(request_data["rely"].(string))
	task_update_data["subtasks_id"] = request_data["subtasks_id"].(string)
	task_update_data["cron"] = request_data["cron"].(string)
	task_update_data["task_exec_type"] = request_data["task_exec_type"].(string)
	task_update_data["exec_task_node_id"] = request_data["exec_task_node_id"].(string)
	task_update_data["cmd"] = request_data["cmd"].(string)
	task_update_data["http_type"] = request_data["http_type"].(string)
	task_update_data["task_fail_num"] = int(request_data["task_fail_num"].(float64))
	task_update_data["task_fail_retry_time"] = int(request_data["task_fail_retry_time"].(float64))
	task_update_data["task_notice"], _ = strconv.Atoi(request_data["task_notice"].(string))
	task_update_data["notice_type"], _ = strconv.Atoi(request_data["notice_type"].(string))
	task_update_data["keyword_notice"] = request_data["keyword_notice"].(string)
	task_update_data["remake"] = request_data["remake"].(string)
	task_update_data["next_exec_time"] = cron_obj.Next(time.Now())

	if task_update_data["exec_task_node_id"] == "0" {
		task_update_data["exec_task_node_type"] = 1
	} else {
		task_update_data["exec_task_node_type"] = 0
	}

	if !task_model.Edit(c.Orm, task_update_data) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	// 获取任务状态
	task_data, _ := task_model.TaskIdToData(c.Orm, task_id)
	if task_data.Status == 1 {
		// 判断有没有在执行的，如果有在执行的，则让正在执行的执行完成，下次执行则执行最新的任务配置
		c.addTaskEtcd(task_data, true)
	}

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG, request_data)
}

// 删除
func (c *TaskController) Del() {
	request_data := make(map[string]float64)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["task_id"] == 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	task_model := models.NewTask()
	task_model.Id = int64(request_data["task_id"])
	if !task_model.Edit(c.Orm, map[string]interface{}{"status": -1}) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.taskEtcd(task_constant.TASK_STOP_DIR+strconv.Itoa(int(request_data["task_id"])), task_model.Id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 启动或停止任务
func (c *TaskController) RunOrStop() {
	request_data := make(map[string]float64)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["task_id"] == 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}
	if !((request_data["status"] == 1) || (request_data["status"] == 2)) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	task_model := models.NewTask()
	task_model.Id = int64(request_data["task_id"])
	c.Orm.Begin()

	if !task_model.Edit(c.Orm, map[string]interface{}{"status": request_data["status"]}) {
		c.Orm.Rollback()
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	if request_data["status"] == 2 {
		c.taskEtcd(task_constant.TASK_STOP_DIR+strconv.Itoa(int(task_model.Id)), task_model.Id)
	} else if request_data ["status"] == 1 {
		task_data, err := task_model.TaskIdToData(c.Orm, task_model.Id)
		if err != nil {
			c.Orm.Rollback()
			c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
		}

		c.addTaskEtcd(task_data, false)
	}

	c.Orm.Commit()
	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 手动执行
func (c *TaskController) ManualExec() {
	// TODO 节点实现了后在回来实现手动执行功能
}

// 强杀任务，强制停止任务执行
func (c *TaskController) Kill() {
	request_data := make(map[string]float64)
	json.Unmarshal(c.Ctx.Input.RequestBody, &request_data)

	if request_data["task_id"] == 0 {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	task_model := models.NewTask()
	task_model.Id = int64(request_data["task_id"])
	if !task_model.Edit(c.Orm, map[string]interface{}{"status": 2}) {
		c.ResponseJson(helpers.FAIL_CODE, false, helpers.FAIL_MSG)
	}

	c.taskEtcd(task_constant.TASK_KILL_DIR+strconv.Itoa(int(request_data["task_id"])), task_model.Id)

	c.ResponseJson(helpers.SUCCESS_CODE, true, helpers.SUCCESS_MSG)
}

// 任务日志
func (c *TaskController) TaskLogList() {

}

// etcd添加任务
func (c *TaskController) addTaskEtcd(task *models.Task, is_edit bool) bool {
	key := task_constant.TASK_LIST_DIR + strconv.Itoa(int(task.Id))

	if is_edit {
		task_stop_key := task_constant.TASK_STOP_DIR + strconv.Itoa(int(task.Id))
		c.taskEtcd(task_stop_key, task.Id)
	}

	c.taskEtcd(key, task.Id)

	return true
}

func (c *TaskController) taskEtcd(key string, task_id int64) {
	task_model := models.NewTask()
	task_data, err := task_model.TaskIdToData(c.Orm, task_id)

	// 任务数据转map格式
	task_data_map := make(map[string]interface{})
	task_data_json, _ := json.Marshal(task_data)
	json.Unmarshal(task_data_json, &task_data_map)

	task_node_model := models.NewTaskNode()

	// 查找节点数据
	if task_data.ExecTaskNodeId != "0" {
		exec_task_node_id_arr := strings.Split(task_data.ExecTaskNodeId, ",")
		task_node_data := task_node_model.IdInData(c.Orm, exec_task_node_id_arr)
		task_data_map["exec_task_node_data"] = task_node_data
	}

	// 查找子任务数据
	if task_data.TaskType == 1 && task_data.SubtasksId != "" {
		subtasks_id_arr := strings.Split(task_data.SubtasksId, ",")
		task_subtasks_data := task_model.IdInData(c.Orm, subtasks_id_arr)

		//获取子任务节点数据
		for k, v := range task_subtasks_data {
			if v["ExecTaskNodeId"].(string) != "0" {
				exec_task_node_id_arr := strings.Split(v["ExecTaskNodeId"].(string), ",")
				task_node_data := task_node_model.IdInData(c.Orm, exec_task_node_id_arr)
				task_subtasks_data[k]["exec_task_node_data"] = task_node_data
			}
		}

		task_data_map["subtasks_data"] = task_subtasks_data
	}

	if err != nil {
		goto JAMP
	}

	// 数据写入
	if _, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), key, string(task_data_json)); err != nil {
		beego.Error("数据写入失败。", err)
	}

	/*go func() {
		task_data_json, _ := json.Marshal(task_data_map)

		// 创建租约
		lease_grant_resp, err := library.G_conf_etcd_client.Lease.Grant(context.TODO(), 1)
		if err != nil {
			beego.Error("创建租约失败。", err)
		}

		// 自动续租
		keep_alive_chan, err := library.G_conf_etcd_client.Lease.KeepAlive(context.TODO(), lease_grant_resp.ID)
		if err != nil {
			beego.Error("自动续租失败。", err)
		}

		// 数据写入
		if _, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), key, string(task_data_json), clientv3.WithLease(lease_grant_resp.ID)); err != nil {
			beego.Error("数据写入失败。", err)
		}

		go func() {
			// 获取任务数据
			for true {
				if task_data.ExecTaskNodeId == "" {
					library.G_conf_etcd_client.Lease.Revoke(context.TODO(), lease_grant_resp.ID)
					break
				}

				resp, err := library.G_conf_etcd_client.Kv.Get(context.TODO(), key)
				if err == nil {
					for _, v := range resp.Kvs {
						json.Unmarshal(v.Value, &task_data_map)
					}
				} else {
					break
				}

				time.Sleep(3 * time.Second)
			}
		}()

		// 续租回应
		for {
			select {
			case keep_alive_resq := <-keep_alive_chan:
				if keep_alive_resq == nil {
					break
				}
			}
		}
	}()*/

JAMP:
}
