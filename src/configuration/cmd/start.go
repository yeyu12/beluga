package cmd

import (
	beluga_drive "beluga/src/beluga/drive"
	"beluga/src/beluga/helpers"
	"beluga/src/configuration/system"
	"github.com/urfave/cli"
	"runtime"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "启动服务",
	Description: "中心配置服务启动",
	Action:      start,
	Flags: []cli.Flag{
		helpers.StringFlag("configDir, c", "config/", "重定向配置文件路径"),
		helpers.StringFlag("host", "0.0.0.0", "监听地址"),
		helpers.StringFlag("port, p", "9411", "监听端口"),
	},
}

var (
	host string
	port string
)

func start(c *cli.Context) {
	if c.IsSet("configDir") {
		beluga_drive.CONFIG_DIR = c.String("configDir")
		beluga_drive.CONFIG_FILENAME = "configuration_node.ini"
	}
	if c.IsSet("host") {
		host = c.String("host")
	}
	if c.IsSet("port") {
		port = c.String("port")
	}

	run()
}

func run() {
	// 初始化线程
	helpers.InitThreadNum(runtime.NumCPU())

	ch := make(chan bool)

	go system.Run(host, port)

	<-ch
}
