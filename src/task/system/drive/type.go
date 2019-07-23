package drive

import (
	"context"
	"github.com/gorhill/cronexpr"
	"time"
)

// 任务
type Task struct {
	Id                int64      `json:"id"`
	Name              string     `json:"name"`
	CreateTime        time.Time  `json:"create_time"`
	StartTime         time.Time  `json:"start_time"`
	ConsumeTime       float64    `json:"consume_time"`
	Overtime          int        `json:"overtime"`
	LastExecType      int        `json:"last_exec_type"`
	TaskType          int        `json:"task_type"`
	Rely              int        `json:"rely"`
	SubtasksId        string     `json:"subtasks_id"`
	Cron              string     `json:"cron"`
	TaskExecType      string     `json:"task_exec_type"`
	ExecTaskNodeType  int        `json:"exec_task_node_type"`
	ExecTaskNodeId    string     `json:"exec_task_node_id"`
	Cmd               string     `json:"cmd"`
	HttpType          string     `json:"http_type"`
	TaskFailNum       int        `json:"task_fail_num"`
	TaskFailRetryTime int        `json:"task_fail_retry_time"`
	TaskNotice        int        `json:"task_notice"`
	NoticeType        int        `json:"notice_type"`
	KeywordNotice     string     `json:"keyword_notice"`      // 提醒关键字
	Remake            string     `json:"remake"`              // 备注
	Status            int        `json:"status"`              // 状态
	NextExecTime      time.Time  `json:"next_exec_time"`      // 下次执行时间
	SubtasksData      []Task     `json:"subtasks_data"`       // 子任务数据
	ExecTaskNodeData  []TaskNode `json:"exec_task_node_data"` // 执行节点数据
}

// 节点
type TaskNode struct {
	Id         int64     `json:"id"`
	Ip         string    `json:"ip"`
	CreateTime time.Time `json:"create_time"`
	Remake     string    `json:"remake"`
	IsDelete   int       `json:"is_delete"` // 0删除，1正常
}

// 任务调度
type TaskSchedulePlan struct {
	Task     *Task
	Expr     *cronexpr.Expression // cron对象
	NextTime time.Time            // 下次执行时间
	TaskLock *TaskLock
}

// 任务执行状态
type TaskExecuteInfo struct {
	Task       *Task
	TheoryTime time.Time // 理论调度时间
	ActualTime time.Time // 实际调度时间
	CancelCtx  context.Context
	CancelFunc context.CancelFunc //  用于取消执行的cancel函数
}

// 任务执行结果
type TaskExecuteResult struct {
	ExecuteInfo *TaskExecuteInfo
	Output      []byte    // 输出
	Err         error     // 错误信息
	StartTime   time.Time // 开始时间
	EndTime     time.Time // 结束时间
}

// 任务执行日志
type TaskLog struct {
	TaskId       int64     `json:"task_id"`        // 任务ID
	TaskName     string    `json:"task_name"`      // 任务名
	Cmd          string    `json:"cmd"`            // 执行命令
	Output       string    `json:"output"`         // 执行输出信息
	CreateTime   time.Time `json:"create_time"`    // 创建时间，也就是开始执行时间
	EndTime      time.Time `json:"end_time"`       // 执行结束时间
	Err          string    `json:"err"`            // 错误信息
	ConsumeTime  float32   `json:"consume_time"`   // 执行消耗时间，s为单位
	TaskExecType int       `json:"task_exec_type"` // 执行状态，1成功，0失败
	NodeIp       string    `json:"node_ip"`        // 节点IP
}

// 待写入的日志数据
type TaskLogWaitBucket struct {
	Logs []interface{}
}

// 结果通知
//type TaskNotice struct {
//	TaskId int64
//}
