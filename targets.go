package movielite

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/unicode/norm"
)

type Target struct {
	Name string `json:"name"`
}

func (s *Service) getTargets(c *gin.Context) {
	targets, err := s.getDirectories()
	if err != nil {
		log.Error(err)
		content := gin.H{"error": err.Error()}
		c.JSON(http.StatusInternalServerError, content)
		return
	}

	c.JSON(http.StatusOK, targets)
}

func (s *Service) getDirectories() ([]Target, error) {
	target := s.Config.TargetDir
	log.Debugf("target dir %s", target)
	var dirs = []Target{}
	//var err = errors.New("")
	if target != "" {
		files, err := ioutil.ReadDir(target)
		if err != nil {
			log.Error(err)
			return dirs, err
		}
		for _, f := range files {
			if f.IsDir() {
				var dir = Target{}
				dir.Name = toUtf8Nfc(f.Name())
				dirs = append(dirs, dir)
			}
		}
	}
	return dirs, nil
}

func toUtf8Nfc(s string) string {
	res := string(norm.NFC.Bytes([]byte(s)))
	return res
}
