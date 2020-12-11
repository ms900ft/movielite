package main

import (
	"os"

	"github.com/prometheus/common/log"
	"github.com/urfave/cli"

	"ms/movielight/commands"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		commands.StartCommand,
		commands.IndexCommand,
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
