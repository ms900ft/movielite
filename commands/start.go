package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ms900ft/movielite"
	"github.com/ms900ft/movielite/pkg/token"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

// StartCommand is used to register the start cli command

var StartCommand = cli.Command{
	Name:  "start",
	Usage: "Starts web server",
	//Flags:  startFlags,
	Action: startAction,
}

var (
	g errgroup.Group
)

// startAction start the web server and initializes the daemon
func startAction(ctx *cli.Context) error {
	conf := movielite.GetConfig(ctx.GlobalString("config"))
	a := movielite.Service{Config: conf}
	w := movielite.Walker{Config: conf}

	a.Initialize()
	token, err := token.AdminToken(conf.Secret)
	if err != nil {
		log.Fatal(err)
	}
	w.Token = token

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
