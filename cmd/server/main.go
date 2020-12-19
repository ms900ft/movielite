package main

import (
	"os"

	"github.com/prometheus/common/log"
	"github.com/urfave/cli"

	"github.com/ms900ft/movielite/commands"
)

var configPath string

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		commands.StartCommand,
		commands.IndexCommand,
		commands.ScanCommand,
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config file",
			//Destination: &configPath,
			Required: false,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
