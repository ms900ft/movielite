package commands

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ms900ft/movielite"
	"github.com/ms900ft/movielite/models"
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

	expiresAt := time.Now().Add(time.Minute * 1000000).Unix()
	tk := &models.Token{
		Name: "admin",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	if conf.Secret == "" {
		log.Fatal("no secret found")
	}
	tokenString, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		log.Fatal(err)
	}
	w.Token = tokenString

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
