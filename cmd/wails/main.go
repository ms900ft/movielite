package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/ms900ft/movielite"
	"github.com/ms900ft/movielite/pkg/token"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend
var frontend embed.FS

//go:embed icon.png
var iconData []byte

var (
	appWindow    *application.WebviewWindow
	movieliteSvc *movielite.Service
	walker       *movielite.Walker
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting Movielite desktop...")

	conf := movielite.GetConfig("")
	conf.Mode = "prod"

	movieliteSvc = &movielite.Service{Config: conf}
	movieliteSvc.Initialize()

	walker = &movielite.Walker{Config: conf}
	tokenStr, err := token.AdminToken(conf.Secret)
	if err != nil {
		log.Fatal(err)
	}
	walker.Token = tokenStr

	go func() {
		if err := movieliteSvc.Run(); err != nil {
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
		movieliteSvc.WorkerPool.Stop()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	subFS, _ := fs.Sub(frontend, "frontend")

	app := application.New(application.Options{
		Name:        "Movielite",
		Description: "Personal Movie Database",
		Assets: application.AssetOptions{
			Handler: http.FileServer(http.FS(subFS)),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	setupSystemTray(app)
	setupMenu(app)

	appWindow = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Movielite",
		Width:  1280,
		Height: 800,
		URL:    "http://localhost:8001/movie2/",
	})
	appWindow.Center()
	appWindow.Show()
	appWindow.Focus()

	setupDragDrop(appWindow)

	app.OnShutdown(func() {
		movieliteSvc.WorkerPool.Stop()
	})

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func setupSystemTray(app *application.App) {
	systray := app.SystemTray.New()

	log.Printf("Icon data length: %d bytes", len(iconData))

	systray.SetIcon(iconData)
	systray.SetLabel("Movielite")
	log.Println("System tray created with icon")

	systray.OnClick(func() {
		if appWindow.IsVisible() {
			appWindow.Hide()
		} else {
			appWindow.Show()
			appWindow.Focus()
		}
	})

	trayMenu := application.NewMenu()
	trayMenu.Add("Show").OnClick(func(ctx *application.Context) {
		appWindow.Show()
		appWindow.Focus()
	})
	trayMenu.Add("Hide").OnClick(func(ctx *application.Context) {
		appWindow.Hide()
	})
	trayMenu.AddSeparator()
	trayMenu.Add("Scan All Folders").OnClick(func(ctx *application.Context) {
		triggerScan("")
	})
	trayMenu.AddSeparator()
	trayMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(trayMenu)
}

func setupMenu(app *application.App) {
	appMenu := app.NewMenu()

	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.Add("Add Folder...").OnClick(func(ctx *application.Context) {
		openFolderDialog(app)
	})
	fileMenu.Add("Scan All").OnClick(func(ctx *application.Context) {
		triggerScan("")
	})
	fileMenu.AddSeparator()
	fileMenu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	editMenu := appMenu.AddSubmenu("Edit")
	editMenu.Add("Copy").OnClick(func(ctx *application.Context) {
		appWindow.ExecJS("document.execCommand('copy')")
	})
	editMenu.Add("Paste").OnClick(func(ctx *application.Context) {
		appWindow.ExecJS("document.execCommand('paste')")
	})
	editMenu.Add("Select All").OnClick(func(ctx *application.Context) {
		appWindow.ExecJS("document.execCommand('selectAll')")
	})

	viewMenu := appMenu.AddSubmenu("View")
	viewMenu.Add("Reload").OnClick(func(ctx *application.Context) {
		appWindow.Reload()
	})
	viewMenu.Add("Toggle DevTools").OnClick(func(ctx *application.Context) {
		appWindow.OpenDevTools()
	})

	settingsMenu := appMenu.AddSubmenu("Settings")
	settingsMenu.Add("Open Config File").OnClick(func(ctx *application.Context) {
		openConfigFile()
	})
	settingsMenu.Add("Reload Config").OnClick(func(ctx *application.Context) {
		reloadConfig()
	})

	app.Menu.Set(appMenu)
}

func setupDragDrop(win *application.WebviewWindow) {
	win.HandleDragEnter()
	win.HandleDragOver(0, 0)
	win.HandleDragLeave()
}

func triggerScan(path string) {
	go func() {
		if path == "" {
			log.Println("Starting full scan...")
			for _, dir := range movieliteSvc.Config.ScanDirectories {
				walker.Run(dir)
			}
		} else {
			log.Printf("Scanning folder: %s", path)
			walker.Run(path)
		}
	}()
}

func openFolderDialog(app *application.App) {
	path, err := app.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		Title:                "Select Movie Folder",
		CanChooseDirectories: true,
		CanChooseFiles:       false,
	}).PromptForSingleSelection()

	if err == nil && path != "" {
		log.Printf("Selected folder: %s", path)
		triggerScan(path)
	}
}

func notifyFrontend(event string, data string) {
	if appWindow != nil {
		apiURL := "http://localhost:8001/api/scan"
		if data != "" {
			apiURL += "?path=" + url.QueryEscape(data)
		}
		appWindow.ExecJS(`window.dispatchEvent(new CustomEvent('` + event + `', {detail: '` + data + `'}))`)
		_, _ = http.Get(apiURL)
	}
}
