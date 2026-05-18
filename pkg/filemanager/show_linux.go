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

func detectFileManager() (string, bool) {
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))

	managerMap := map[string]string{
		"gnome": "gio",
		"unity": "gio",
		"ubuntu": "gio",
		"xfce4": "exo-open",
		"mate": "gio",
		"cinnamon": "gio",
		"kde": "dolphin",
		"plasma": "dolphin",
		"kxfce": "exo-open",
	}

	for _, desk := range []string{"gnome", "unity", "ubuntu", "xfce4", "mate", "cinnamon", "kde", "plasma", "kxfce"} {
		if strings.Contains(desktop, desk) {
			if manager, ok := managerMap[desk]; ok {
				return manager, true
			}
		}
	}

	return "xdg-open", false
}

	for _, desk := range []string{"gnome", "unity", "ubuntu", "xfce4", "mate", "cinnamon", "kde", "plasma", "kxfce"} {
		if strings.Contains(desktop, desk) {
			if manager, ok := managerMap[desk]; ok {
				return manager, true
			}
		}
	}

	return "xdg-open", false
}

func showViaCMD(args ...string) error {
	fullCmd := exec.Command(args[0], args[1:]...)
	return fullCmd.Start()
}

func showFile(path string) error {
	manager, _ := detectFileManager()

	switch manager {
	case "gio":
		return showViaCMD("gio", "open", path)
	case "dolphin":
		return showViaCMD("dolphin", "--select", path)
	case "nautilus":
		return showViaCMD("nautilus", "--select", path)
	case "thunar":
		return showViaCMD("thunar", "--select", path)
	case "konqueror":
		return showViaCMD("konqueror", "--select", path)
	case "spacefm":
		return showViaCMD("spacefm", "--select", path)
	case "pcmanfm":
		return showViaCMD("pcmanfm", "--select", path)
	case "exo-open":
		return showViaCMD("exo-open", "--launch", "FileManager", path)
	default:
		return showViaCMD("xdg-open", path)
	}
}

func ShowDir(path string) error {
	dir := filepath.Dir(path)

	if path != dir {
		return showFile(path)
	}

	manager, _ := detectFileManager()

	var err error
	switch manager {
	case "gio":
		err = showViaCMD("gio", "open", dir)
	case "dolphin":
		err = showViaCMD("dolphin", dir)
	case "nautilus":
		err = showViaCMD("nautilus", dir)
	case "thunar":
		err = showViaCMD("thunar", dir)
	case "konqueror":
		err = showViaCMD("konqueror", dir)
	case "spacefm":
		err = showViaCMD("spacefm", dir)
	case "pcmanfm":
		err = showViaCMD("pcmanfm", dir)
	case "exo-open":
		err = showViaCMD("exo-open", "--launch", "FileManager", dir)
	default:
		err = showViaCMD("xdg-open", dir)
	}

	if err != nil {
		log.Errorf("failed to show directory %s: %v", dir, err)
		return err
	}

	return nil
}

	for _, desk := range []string{"gnome", "unity", "ubuntu", "xfce4", "mate", "cinnamon", "kde", "plasma", "kxfce"} {
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

func ShowDir(dir string) error {
	manager := detectFileManager()

	var err error

	switch manager {
	case "gio":
		err = showViaCMD("gio", "open", dir)
	case "exo-open":
		err = showViaCMD("exo-open", "--launch", "FileManager", dir)
	default:
		err = showViaCMD("xdg-open", dir)
	}

	if err != nil {
		log.Errorf("failed to show directory %s: %v", dir, err)
		return err
	}

	return nil
}
