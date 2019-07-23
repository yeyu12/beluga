package drive

import (
	"beluga/src/beluga/drive"
	"beluga/src/beluga/helpers"
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// 任务调度
type Scheduler struct {
	TaskTable          map[string]*TaskSchedulePlan // 任务调度表,key为Id+_+name
	TaskExecTable      map[string]*TaskExecuteInfo  // 任务执行表
	TaskStopChan       chan *Task                   // 任务停止队列
	TaskDelChan        chan *Task                   // 任务删除队列
	TaskKillChan       chan *Task                   // 任务强杀队列
	//TaskExecManualChan chan *Task                   // 任务手动执队列
	TaskResultChan     chan *TaskExecuteResult      // 任务结果队列
	//TaskNoticeTable chan *task_constant.Task// 任务通知队列
}

var G_scheduler *Scheduler

// 初始化
func InitScheduler() {
	G_scheduler = &Scheduler{
		TaskTable:          make(map[string]*TaskSchedulePlan),
		TaskExecTable:      make(map[string]*TaskExecuteInfo),
		TaskStopChan:       make(chan *Task, 1000),
		TaskDelChan:        make(chan *Task, 1000),
		TaskKillChan:       make(chan *Task, 1000),
		//TaskExecManualChan: make(chan *Task, 1000),
		TaskResultChan:     make(chan *TaskExecuteResult, 1000),
	}

	// 任务调度
	go G_scheduler.scheduleLoop()
}

// 任务执行结果写入到结果队列中
func (scheduler *Scheduler) PushTaskResult(task_result *TaskExecuteResult) {
	scheduler.TaskResultChan <- task_result
}

// 组装任务key
func (scheduler *Scheduler) AssTaskKey(task *Task) string {
	return strconv.Itoa(int(task.Id)) + "_" + task.Name
}

// 任务调度
func (scheduler *Scheduler) scheduleLoop() {
	schedule_after_func := func() (scheduleAfter time.Duration) {
		var (
			near_time *time.Time
			time_now  = time.Now() // 当前时间
		)

		if len(G_scheduler.TaskTable) == 0 {
			scheduleAfter = 1 * time.Second
			return
		}
		for _, v := range G_scheduler.TaskTable {
			if v.NextTime.Before(time_now) || v.NextTime.Equal(time_now) {
				// 执行调用
				scheduler.execTaskStartInit(v)

				// 更新下次时间
				v.NextTime = v.Expr.Next(time_now)
			}

			// 获取最近快要过期的时间，或已经过期的时间
			if near_time == nil || v.NextTime.Before(*near_time) {
				near_time = &v.NextTime
			}
		}

		scheduleAfter = (*near_time).Sub(time_now)

		return
	}

	scheduler_timer := time.NewTimer(schedule_after_func())

	for {
		select {
		case task_result := <-scheduler.TaskResultChan:
			scheduler.handleTaskResult(task_result)
		case task_stop := <-scheduler.TaskStopChan:
			scheduler.handleStopTask(task_stop)
		case task_kill := <-scheduler.TaskKillChan:
			scheduler.handleKillTask(task_kill)
		case <-scheduler_timer.C:
		}

		scheduler_timer.Reset(schedule_after_func())
	}
}

// 任务执行初始化
func (scheduler *Scheduler) execTaskStartInit(task_plan *TaskSchedulePlan) {
	task_exec_info := &TaskExecuteInfo{
		Task:       task_plan.Task,
		TheoryTime: task_plan.NextTime,
		ActualTime: time.Now(),
	}

	task_exec_info.CancelCtx, task_exec_info.CancelFunc = context.WithTimeout(context.TODO(), 2*time.Second)
	G_scheduler.TaskExecTable[scheduler.AssTaskKey(task_plan.Task)] = task_exec_info

	// 调用任务执行
	G_executor.ExecTask(task_exec_info)
}

// 处理任务结果
func (scheduler *Scheduler) handleTaskResult(task_result *TaskExecuteResult) {
	delete(scheduler.TaskExecTable, scheduler.AssTaskKey(task_result.ExecuteInfo.Task))

	local_ip, err := helpers.GetLocalIp()
	if err != nil {
		drive.Notices(logrus.Fields{}, errors.Wrapf(err, "获取节点IP失败"))
	}
	// 生成日志数据
	log_data := &TaskLog{
		TaskId:       task_result.ExecuteInfo.Task.Id,
		TaskName:     task_result.ExecuteInfo.Task.Name,
		Cmd:          task_result.ExecuteInfo.Task.Cmd,
		Output:       string(task_result.Output),
		CreateTime:   task_result.StartTime,
		EndTime:      task_result.EndTime,
		Err:          "",
		ConsumeTime:  float32((task_result.EndTime.UnixNano() - task_result.StartTime.UnixNano())) / 1e9,
		TaskExecType: 1,
		NodeIp:       local_ip,
	}

	if task_result.Err != nil {
		log_data.Err = task_result.Err.Error()
	}

	if task_result.Err != nil {
		log_data.TaskExecType = 0
	}

	G_task_log_sink.PushLog(log_data)
}

//任务停止
func (scheduler *Scheduler) handleStopTask(task_data *Task) {
	key := scheduler.AssTaskKey(task_data)

	if task_info, task := scheduler.TaskTable[key]; task {
		task_info.TaskLock.Unlock()
	}

	delete(scheduler.TaskTable, key) // 删除键
}

//任务强杀
func (scheduler *Scheduler) handleKillTask(task_data *Task) {
	key := scheduler.AssTaskKey(task_data)

	if task_exec_info, task_executing := scheduler.TaskExecTable[key]; task_executing {
		task_exec_info.CancelFunc()
	}

	if task_info, task := scheduler.TaskTable[key]; task {
		task_info.TaskLock.Unlock()
	}

	delete(scheduler.TaskTable, key) // 删除键
}
