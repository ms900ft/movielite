package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"

	"ms/movielight"
)

// StartCommand is used to register the start cli command
var StartCommand = cli.Command{
	Name:  "start",
	Usage: "Starts web server",
	//Flags:   startFlags,
	Action: startAction,
}

var startFlags = []cli.Flag{
	// cli.BoolFlag{
	// 	Name:  "detach-server, d",
	// 	Usage: "detach from the console (daemon mode)",
	//	},
}
var (
	g errgroup.Group
)

// startAction start the web server and initializes the daemon
func startAction(ctx *cli.Context) error {

	conf := movielight.GetConfig()
	a := movielight.Service{Config: conf}
	w := movielight.Walker{Config: conf}

	a.Initialize()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		a.WorkerPool.Stop()
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()
	//a.Run(":8080")
	g.Go(func() error {
		return a.Run()
	})

	g.Go(func() error {
		return w.RunWatcher()
	})
	// if viper.GetBool("Rescan.Enable") {
	// 	g.Go(func() error {
	// 		return r.Run()
	// 	})
	// }
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	return nil
}
