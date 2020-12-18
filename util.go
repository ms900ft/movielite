package movielite

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func URLEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

const raw = `
on run argv
  tell application "Finder"
    repeat with f in argv
      move (f as POSIX file) to trash
    end repeat
  end tell
end run
`

//var trash = filepath.Join("/Volumes/Transcend/.Trashes", "501")

func Trash(f string, trash string) (trashcan string, err error) {
	bin, err := exec.LookPath("osascript")
	if err != nil {
		err = fmt.Errorf("not yet supported")
		return
	}

	if _, err = os.Stat(trash); err != nil {
		err = fmt.Errorf("trash not found")
		return
	}

	path, err := filepath.Abs(f)
	if err != nil {
		return
	}
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	_ = name

	dest := filepath.Join(trash, base)
	if _, err = os.Stat(dest); err == nil {
		err = fmt.Errorf("already exists")
		return
	}
	trashcan = dest
	log.Debug(path)
	params := append([]string{"-e", raw}, path)
	cmd := exec.Command(bin, params...)
	log.Debugf("%+v", params)
	if err = cmd.Run(); err != nil {
		log.Error(err)
		return
	}

	if _, err = os.Stat(trashcan); err != nil {
		trashcan = ""
	}

	return trashcan, err
}
