package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ms900ft/movielite"
	"github.com/ms900ft/movielite/pkg/token"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend
var frontend embed.FS

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting Movielite desktop...")

	conf := movielite.GetConfig("")
	conf.Mode = "prod"

	svc := &movielite.Service{Config: conf}
	svc.Initialize()

	walker := &movielite.Walker{Config: conf}
	tokenStr, err := token.AdminToken(conf.Secret)
	if err != nil {
		log.Fatal(err)
	}
	walker.Token = tokenStr

	go func() {
		if err := svc.Run(); err != nil {
			log.Println("Server error:", err)
		}
	}()

	go func() {
		if err := walker.RunWatcher(); err != nil {
			log.Println("Watcher error:", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signalChan
		svc.WorkerPool.Stop()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	log.Printf("Movielite desktop started on http://localhost:%d", conf.Port)
	log.Printf("Window should open automatically. Access at http://localhost:%d/movie2/", conf.Port)

	// Get frontend fs
	subFS, _ := fs.Sub(frontend, "frontend")

	app := application.New(application.Options{
		Name:        "Movielite",
		Description: "Personal Movie Database",
		Assets: application.AssetOptions{
			Handler: http.FileServer(http.FS(subFS)),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Movielite",
		Width:  1280,
		Height: 800,
		URL:    "http://localhost:8001/movie2/",
	})
	window.Center()
	window.Show()
	window.Focus()

	app.OnShutdown(func() {
		svc.WorkerPool.Stop()
	})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
