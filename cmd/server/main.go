package main

import (
	"fmt"
	"log"
	"ms/movielight"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	a := movielight.Service{}
	//w := movielight.Walker{}
	//r := movielight.Resanner{}
	viper.SetConfigName("movielight")
	//viper.AddConfigPath("./cmd/films")
	//viper.AddConfigPath("/Users/ms/moviedb")
	viper.AddConfigPath(".")
	//viper.AddConfigPath("/etc/moviedbconfig")
	//viper.AddConfigPath("/Users/ms")
	viper.SetConfigType("yaml")
	viper.SetDefault("MovieServerUrl", "http://localhost:8000")
	viper.SetDefault("Rescan.Delay", 10)
	viper.SetDefault("language", "de-DE")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("error config file: %s", err))
	}

	db := viper.GetString("Database.Dbname")

	a.Initialize(
		db)

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
		return a.Run(":8080")
	})

	// g.Go(func() error {
	// 	return w.Run()
	// })
	// if viper.GetBool("Rescan.Enable") {
	// 	g.Go(func() error {
	// 		return r.Run()
	// 	})
	// }
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
