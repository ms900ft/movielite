//go:build linux
// +build linux

package filemanager

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func detectFileManager() string {
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))

	managerMap := map[string]string{
		"gnome":      "gio",
		"unity":      "gio",
		"ubuntu":     "gio",
		"xfce4":      "exo-open",
		"mate":       "gio",
		"cinnamon":   "gio",
		"kde":        "gio",
		"plasma":     "gio",
		"kxfce":      "exo-open",
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
