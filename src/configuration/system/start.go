package system

import (
	"beluga/src/beluga/configuration_constant"
	beluga_drive "beluga/src/beluga/drive"
	"beluga/src/beluga/helpers"
	"beluga/src/beluga/library"
	"beluga/src/configuration/system/drive"
)

// 初始化各种服务
func initService() {
	beluga_drive.G_node_conf = make(map[string]string)

	if beluga_drive.CONFIG_DIR == "" && beluga_drive.CONFIG_FILENAME == "" {
		beluga_drive.CONFIG_DIR = helpers.GetCurrentDirectory() + "/../config/"
		beluga_drive.CONFIG_FILENAME = "configuration_node.ini"
	}

	// 配置
	beluga_drive.InitConfig()

	// 日志初始化
	beluga_drive.InitLog(helpers.GetCurrentDirectory() + "/../")

	// redis
	beluga_drive.InitRedis()

	// 系统信息
	beluga_drive.InitMonitor()

	// 数据库
	beluga_drive.InitMysql()

	// etcd
	beluga_drive.InitEtcd()

	// 服务注册
	go library.ServerRegister(configuration_constant.CONFIGURATION_REGISTER_DIR, "")

	// 配置初始化
	drive.InitNodeConf()
}

func Run(host, port string) {
	initService()

	// 监听配置
	drive.WatchNode()
	drive.WatchConfigurationReleaseToSync()
	drive.WatchNodeConf()

	// 调用http
	HttpServer(host, port)
}
