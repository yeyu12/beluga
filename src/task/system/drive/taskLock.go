package drive

import (
	"beluga/src/beluga/task_constant"
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
)

// 任务锁
type TaskLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	taskId     string             // 任务ID
	cancelFunc context.CancelFunc // 用于终止自动续租
	leaseId    clientv3.LeaseID   // 租约ID
	isLocked   bool               // 是否上锁成功
}

// 初始化一把锁
func InitTaskLock(task_id string, kv clientv3.KV, lease clientv3.Lease) (jobLock *TaskLock) {
	jobLock = &TaskLock{
		kv:     kv,
		lease:  lease,
		taskId: task_id,
	}

	return
}

// 尝试上锁
func (taskLock *TaskLock) TryLock() (err error) {
	var (
		txn      clientv3.Txn
		lock_key string
		txn_resp *clientv3.TxnResponse
	)

	// 创建租约
	lease_resp, err := taskLock.lease.Grant(context.TODO(), 5)
	if err != nil {
		return
	}

	cancel_ctx, cancel_func := context.WithCancel(context.TODO())
	keep_resp_chan, err := taskLock.lease.KeepAlive(cancel_ctx, lease_resp.ID)
	if err != nil {
		goto JAMP
	}

	// 续租回应
	go func() {
		var keep_resp *clientv3.LeaseKeepAliveResponse
		for {
			select {
			case keep_resp = <-keep_resp_chan: // 自动续租的应答
				if keep_resp == nil {
					goto KEEP_END
				}
			}
		}

	KEEP_END:
	}()

	// 事务操作
	txn = taskLock.kv.Txn(context.TODO())
	lock_key = task_constant.TASK_LOCK_DIR + taskLock.taskId

	txn.If(clientv3.Compare(clientv3.CreateRevision(lock_key), "=", 0)).
		Then(clientv3.OpPut(lock_key, "", clientv3.WithLease(lease_resp.ID))).
		Else(clientv3.OpGet(lock_key))

	// 提交事务
	if txn_resp, err = txn.Commit(); err != nil {
		goto JAMP
	}

	if !txn_resp.Succeeded {
		err = errors.New("锁被占用")
		goto JAMP
	}

	taskLock.leaseId = lease_resp.ID
	taskLock.cancelFunc = cancel_func
	taskLock.isLocked = true

	return
JAMP:
	cancel_func()
	taskLock.lease.Revoke(context.TODO(), lease_resp.ID)

	return
}

// 释放锁
func (taskLock *TaskLock) Unlock() {
	if taskLock.isLocked {
		taskLock.cancelFunc()
		taskLock.lease.Revoke(context.TODO(), taskLock.leaseId)
	}
}
