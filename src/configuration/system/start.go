package system

import (
	"beluga/src/configuration/system/drive"
	"beluga/src/beluga/library"
	"beluga/src/beluga/configuration_constant"
)

// 初始化各种服务
func initService() {
	drive.G_node_conf = make(map[string]string)

	// 配置
	drive.InitConfig()

	// 初始化日子还
	drive.InitLog()

	// redis
	drive.InitRedis()

	// 系统信息
	drive.InitMonitor()

	// 数据库
	drive.InitMysql()

	// etcd
	drive.InitEtcd()

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