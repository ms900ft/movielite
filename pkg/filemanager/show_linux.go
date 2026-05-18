//go:build linux
// +build linux

package filemanager

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func detectFileManager() string {
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))

	managerMap := map[string]string{
		"gnome":   "gio",
		"unity":   "gio",
		"ubuntu":  "gio",
		"xfce4":   "exo-open",
		"mate":    "gio",
		"cinnamon": "gio",
		"kde":     "xdg-open",
		"plasma":  "xdg-open",
	}

	for _, desk := range []string{"gnome", "unity", "ubuntu", "xfce4", "mate", "cinnamon", "kde", "plasma"} {
		if strings.Contains(desktop, desk) {
			if manager, ok := managerMap[desk]; ok {
				return manager
			}
		}
	}

	return "xdg-open"
}

func showViaCMD(args ...string) error {
	fullCmd := exec.Command(args[0], args[1:]...)
	return fullCmd.Start()
}

func ShowDir(path string) error {
	dir := filepath.Dir(path)

	var err error
	manager := detectFileManager()

	switch manager {
	case "gio":
		err = showViaCMD("gio", "open", path)
	case "nautilus":
		err = showViaCMD("nautilus", path)
	case "thunar":
		err = showViaCMD("thunar", path)
	case "dolphin":
		err = showViaCMD("dolphin", path)
	case "konqueror":
		err = showViaCMD("konqueror", path)
	case "spacefm":
		err = showViaCMD("spacefm", path)
	case "pcmanfm":
		err = showViaCMD("pcmanfm", path)
	case "exo-open":
		err = showViaCMD("exo-open", "--launch", "FileManager", path)
	default:
		err = showViaCMD("xdg-open", path)
	}

	if err != nil {
		log.Errorf("failed to show directory %s: %v", dir, err)
		return err
	}

	return nil
}
