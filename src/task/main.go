package main

import (
	"github.com/urfave/cli"
	"os"
	"beluga/src/task/cmd"
)

const BELUGA_TASK_VERSION = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "task"
	app.Version = BELUGA_TASK_VERSION
	app.Usage = "定时任务"
	app.Commands = []cli.Command{
		cmd.Start,
	}

	app.Run(os.Args)
}