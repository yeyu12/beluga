package main

import (
	"beluga/src/web_server/cmd"
	"github.com/astaxie/beego"
	"github.com/urfave/cli"
	"os"
)

const BELUGA_MASTER_VERSION = "0.0.1"

func main() {
	if beego.AppConfig.DefaultBool("debug", false) {
		cmd.Run()
	} else {
		app := cli.NewApp()
		app.Name = "beluga"
		app.Version = BELUGA_MASTER_VERSION
		app.Usage = "蓝鲸管理后台"
		app.Commands = []cli.Command{
			cmd.Start,
		}

		app.Run(os.Args)
	}
}