package main

import (
	"os"

	"github.com/prometheus/common/log"
	"github.com/urfave/cli"

	"ms/movielite/commands"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		commands.StartCommand,
		commands.IndexCommand,
		commands.ScanCommand,
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
