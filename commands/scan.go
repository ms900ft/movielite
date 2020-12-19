package commands

import (
	"github.com/prometheus/common/log"
	"github.com/urfave/cli"

	"ms/movielite"
)

var directory string

// StartCommand is used to register the start cli command
var ScanCommand = cli.Command{
	Name:   "scan",
	Usage:  "scan directory for movies",
	Flags:  scanFlags,
	Action: scanAction,
}

var scanFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "directory, d",
		Usage:       "scan directory for movies",
		Destination: &directory,
		Required:    true,
	},
}

// startAction start the web server and initializes the daemon
func scanAction(ctx *cli.Context) error {

	conf := movielite.GetConfig(ctx.GlobalString("config"))

	w := movielite.Walker{Config: conf}
	err := w.Run(directory)
	if err != nil {
		log.Fatalf("can't scan for movies: %s", err)
	}
	return nil
}
