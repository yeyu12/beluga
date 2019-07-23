package drive

import (
	"beluga/src/beluga/drive"
	"beluga/src/beluga/library"
	"beluga/src/beluga/task_constant"
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

// 任务管理
type TaskManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var G_task_manager *TaskManager

func InitTaskManager() {
	G_task_manager = &TaskManager{
		client:  library.G_conf_etcd_client.Client,
		kv:      library.G_conf_etcd_client.Kv,
		lease:   library.G_conf_etcd_client.Lease,
		watcher: clientv3.NewWatcher(library.G_conf_etcd_client.Client),
	}

	G_task_manager.WatchTaskList()
	G_task_manager.WatchTaskStop()
	G_task_manager.WatchTaskKill()
	//G_task_manager.WatchTaskManualExec()
}

// 监听任务列表
func (taskManager *TaskManager) WatchTaskList() {
	key := task_constant.TASK_LIST_DIR
	get_resp, err := library.G_conf_etcd_client.Client.Get(context.TODO(), key, clientv3.WithPrefix())

	if err != nil {
		drive.Err(logrus.Fields{}, errors.Wrap(err, "任务列表获取失败。"))
	}

	go func() {
		watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithRev(get_resp.Header.Revision+1), clientv3.WithPrefix())

		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "任务列表监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					taskManager.taskHandle(resp)
					break
				case mvccpb.DELETE:
					break

				}
			}
		}
	}()
}

// 任务列表处理
func (taskManager *TaskManager) taskHandle(event *clientv3.Event) {
	key_arr := strings.Split(string(event.Kv.Key), "/")
	task_id := key_arr[3]
	//local_ip, _ := helpers.GetLocalIp()

	var task_data *Task
	task_data = &Task{}
	json.Unmarshal(event.Kv.Value, &task_data)

	// 获取任务是否在执行
	exec_server_key := task_constant.TASK_EXEC_SERVER_DIR + task_id + "/"

	_, err := library.G_conf_etcd_client.Kv.Get(context.TODO(), exec_server_key, clientv3.WithPrefix())

	if err != nil {
		drive.Err(logrus.Fields{}, errors.Wrap(err, "任务执行列表获取失败"))
		return
	}

	/*
		分为两步
			1、随机节点的话，随机写入。给该任务加锁，排它
			2、指定方式运行。直接判断IP是否存在，如果存在则执行，ip不存在跳出
	*/
	/*
		父任务为随机节点
			子任务为随机节点，在父任务执行所在的节点执行
			子任务为指定节点，下发随机节点，运行的节点数据保存到etcd和父任务数据中
		父任务为指定节点
			子任务为随机节点，在父任务执行所在的节点执行
			子任务为指定节点，下发随机节点，运行的节点数据保存到etcd和父任务数据中
	*/
	if task_data.ExecTaskNodeId == "0" {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		lock_obj := InitTaskLock(task_id, library.G_conf_etcd_client.Kv, library.G_conf_etcd_client.Lease)

		if lock_obj.TryLock() != nil {
			return
		} else {
			// 写入执行目录，数据写入到内存中去
			cron_expr_obj := cronexpr.MustParse(task_data.Cron)

			G_scheduler.TaskTable[G_scheduler.AssTaskKey(task_data)] = &TaskSchedulePlan{
				Task:     task_data,
				Expr:     cron_expr_obj,
				NextTime: cron_expr_obj.Next(time.Now()),
				TaskLock: lock_obj,
			}
		}
	}
}

// 监听任务停止
func (taskManager *TaskManager) WatchTaskStop() {
	key := task_constant.TASK_STOP_DIR
	watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithPrevKV(), clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "停止任务监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					taskManager.taskStop(resp)
					break
				case mvccpb.DELETE:
					break

				}
			}
		}
	}()
}

// 任务停止写入停止队列
func (taskManager *TaskManager) taskStop(event *clientv3.Event) {
	var task_data *Task
	json.Unmarshal(event.Kv.Value, &task_data)

	G_scheduler.TaskStopChan <- task_data
}

// 监听任务强杀
func (taskManager *TaskManager) WatchTaskKill() {
	key := task_constant.TASK_KILL_DIR
	watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithPrevKV(), clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "强杀任务监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					taskManager.taskKill(resp)
					break
				case mvccpb.DELETE:
					break

				}
			}
		}
	}()
}

// 任务强杀写入队列
func (taskManager *TaskManager) taskKill(event *clientv3.Event) {
	var task_data *Task
	json.Unmarshal(event.Kv.Value, &task_data)

	G_scheduler.TaskKillChan <- task_data
}

//监听手动执行
/*func (taskManager *TaskManager) WatchTaskManualExec() {
	key := task_constant.TASK_EXEC_DIR
	watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithPrevKV(), clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "手动执行任务监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					break
				case mvccpb.DELETE:
					break

				}
			}
		}
	}()
}

// 任务手动执行
func (taskManager *TaskManager) taskManualExec(event *clientv3.Event) {

}
*/