package task_constant

var (
	// 配置节点发现目录
	TASK_REGISTER_DIR = "/task/register/"

	// 任务列表目录
	// 任务修改后，节点当前任务执行完成后中断下次执行，把修改后的任务添加进入到列表中来
	TASK_LIST_DIR = "/task/list/" // +id=任务内容

	// 任务执行目录
	TASK_EXEC_SERVER_DIR = "/task/exec/server/" // +id/+ip/=json数组，运行的节点ip

	// 手动执行目录
	TASK_EXEC_DIR = "/task/exec/manual/" // +id/="" 租约的形式下发到节点后，节点删除该执行。

	// 停止任务目录
	TASK_STOP_DIR = "/task/stop/" // +id="" 租约的形式下发到节点后，删除任务

	// 强杀目录
	TASK_KILL_DIR = "/task/kill/" // +id="" 租约的形式下发到节点后，删除

	// 任务锁目录
	TASK_LOCK_DIR = "/task/lock/" // +id="" 随机任务执行
)
