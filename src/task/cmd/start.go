package cmd

import (
	beluga_drive "beluga/src/beluga/drive"
	"beluga/src/beluga/helpers"
	"beluga/src/task/system"
	"github.com/urfave/cli"
	"runtime"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "启动服务",
	Description: "task服务启动",
	Action:      start,
	Flags: []cli.Flag{
		helpers.StringFlag("configDir, c", "config/", "重定向配置文件路径"),
	},
}

func start(c *cli.Context) {
	if c.IsSet("configDir") {
		beluga_drive.CONFIG_DIR = c.String("configDir")
		beluga_drive.CONFIG_FILENAME = "task.ini"
	}

	run()
}

func run() {
	// 初始化线程
	helpers.InitThreadNum(runtime.NumCPU())

	/*go func() {
		for true {
			fmt.Println(runtime.NumGoroutine())
			time.Sleep(1 * time.Second)
		}
	}()*/

	ch := make(chan bool)

	go system.Run()

	<-ch
}
