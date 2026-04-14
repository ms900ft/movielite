//go:build !wails

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/ms900ft/movielite/commands"
)

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
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
