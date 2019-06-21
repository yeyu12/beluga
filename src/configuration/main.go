package main

import (
	"github.com/urfave/cli"
	"os"
	"beluga/src/configuration/cmd"
)

const BELUGA_CONFIGURATION_VERSION = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "configuration"
	app.Version = BELUGA_CONFIGURATION_VERSION
	app.Usage = "中心配置"
	app.Commands = []cli.Command{
		cmd.Start,
	}

	app.Run(os.Args)
}