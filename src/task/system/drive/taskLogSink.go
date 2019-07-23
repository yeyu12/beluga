package drive

import (
	"beluga/src/beluga/drive"
	"beluga/src/web_server/helpers"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

// 任务执行日志
type TaskLogSink struct {
	logsChan       chan *TaskLog           // 日志写入队列
	autoCommitChan chan *TaskLogWaitBucket // 提交
}

var G_task_log_sink *TaskLogSink

// 初始化
func InitTaskLogSink() {
	G_task_log_sink = &TaskLogSink{
		logsChan:       make(chan *TaskLog, 1000),
		autoCommitChan: make(chan *TaskLogWaitBucket, 1000),
	}

	go G_task_log_sink.writeLogPoke()
}

// 日志写入队列
func (taskLogSink *TaskLogSink) PushLog(task_log *TaskLog) {
	select {
	case taskLogSink.logsChan <- task_log:
	default:
	}
}

// 日志写入数据库
func (taskLogSink *TaskLogSink) saveLog(task_log_bucket *TaskLogWaitBucket) {
	var (
		sql         = "INSERT INTO `" + helpers.GetTablePrefix() + "task_log` "
		table_field = ""
	)

	go func() {
		for _, val := range task_log_bucket.Logs {
			task_log := val.(*TaskLog)
			typeofCat := reflect.TypeOf(*task_log)

			if table_field == "" {
				sql += "("
				for i := 0; i < typeofCat.NumField(); i++ {
					field_type := typeofCat.Field(i)
					table_field += "`" + field_type.Tag.Get("json") + "`" + ","
				}
				sql += table_field
				sql = sql[0 : len(sql)-1]
				sql += ") VALUES "
			}

			sql += "("
			value_of := reflect.ValueOf(*task_log)
			for i := 0; i < typeofCat.NumField(); i++ {
				field_type := typeofCat.Field(i)
				sql += "\"" + fmt.Sprint(value_of.FieldByName(field_type.Name)) + "\","
			}
			sql = sql[0 : len(sql)-1]
			sql += "),"
		}
		sql = sql[0 : len(sql)-1]

		_, err := drive.G_mysql.DB().Exec(sql)
		if err != nil {
			drive.Notices(logrus.Fields{}, errors.Wrapf(err, "运行日志写入数据库失败"))
		}
	}()
}

// 日志写入存储捅中
func (taskLogSink *TaskLogSink) writeLogPoke() {
	var (
		log         *TaskLog
		logBucket   *TaskLogWaitBucket
		commitTimer *time.Timer
	)

	for {
		select {
		case log = <-taskLogSink.logsChan:
			if logBucket == nil {
				logBucket = &TaskLogWaitBucket{}

				commitTimer = time.AfterFunc(1*time.Second,
					func(bucket *TaskLogWaitBucket) func() {
						return func() {
							taskLogSink.autoCommitChan <- logBucket
						}
					}(logBucket))
			}

			logBucket.Logs = append(logBucket.Logs, log)

			// 大于或等于数量或超时则自动写入
			if len(logBucket.Logs) >= 3 {
				taskLogSink.saveLog(logBucket)
				commitTimer.Stop()
				logBucket = nil
			}
		case timeout_bucket := <-taskLogSink.autoCommitChan:
			taskLogSink.saveLog(timeout_bucket)
			logBucket = nil
		}

	}
}
