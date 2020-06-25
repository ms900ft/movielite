package main

import (
	"fmt"
	"log"
	"ms/movielight"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	c := movielight.GetConfig()
	a := movielight.Service{Config: c}
	w := movielight.Walker{}

	a.Initialize()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
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
		return w.Run()
	})
	// if viper.GetBool("Rescan.Enable") {
	// 	g.Go(func() error {
	// 		return r.Run()
	// 	})
	// }
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
