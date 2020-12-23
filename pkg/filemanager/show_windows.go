// +build windows

package filemanager

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ShowDir(dir string) error {
	err := exec.Command("explorer", `/select,`, dir).Start()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
