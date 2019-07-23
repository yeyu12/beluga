package system

import (
	beluga_drive "beluga/src/beluga/drive"
	"beluga/src/beluga/helpers"
	"beluga/src/beluga/library"
	"beluga/src/beluga/task_constant"
	"beluga/src/task/system/drive"
	"fmt"
)

// 初始化各种服务
func initService() {
	beluga_drive.G_node_conf = make(map[string]string)

	if beluga_drive.CONFIG_DIR == "" && beluga_drive.CONFIG_FILENAME == "" {
		beluga_drive.CONFIG_DIR = helpers.GetCurrentDirectory() + "/../config/"
		beluga_drive.CONFIG_FILENAME = "task.ini"
	}

	// 配置
	beluga_drive.InitConfig()

	// 日志
	beluga_drive.InitLog(helpers.GetCurrentDirectory() + "/../")

	// redis
	//beluga_drive.InitRedis()

	// 系统信息
	beluga_drive.InitMonitor()

	// 数据库
	beluga_drive.InitMysql()

	// etcd
	beluga_drive.InitEtcd()

	// 服务注册
	go library.ServerRegister(task_constant.TASK_REGISTER_DIR, "")

	// 任务初始化

}

func Run() {
	initService()

	// 执行日志
	drive.InitTaskLogSink()

	drive.InitScheduler()

	drive.InitExecutor()

	// 任务监听
	drive.InitTaskManager()

	// 调用http
	//HttpServer(host, port)

	fmt.Println("服务已启动")
}
