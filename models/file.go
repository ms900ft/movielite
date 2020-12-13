package models

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type File struct {
	gorm.Model
	FullPath string `json:"FullPath" gorm:"Type:text;UNIQUE;index"`
	FileName string `json:"FileName"`
	//CreatedAt time.Time
}

var defaultRegexes = map[string]string{
	"111.OTR": `^(.+?)_\d{2}\.\d{2}\.\d{2}_\d{2}-\d{2}_`,
	"1.OTR":   `^(.+?)_Top_Tipp_+\d{2}\.\d{2}\.\d{2}_\d{2}-\d{2}_`,
	"22.MT":   `^.+?-(.+?)-\d+-\d+.+\.\w{3}$`,
	"21.MT":   `^.+?_-_.+?-(.+)-\d+-\d+.+\.\w{3}$`,
	//	"MZ2":           `^_.+?-(.+)`,
	"1.ARTE_FERNSEHFILM": `.+?-_Fernsehfilme-(.+?)-\d+-\d+.+\.\w{3}$`,
	"3.MT_ZDF":           `^.+?-(.+?)_-`,
	"4.MT_ZDF_SEASON":    `^.+?-(.+?)_\((\d+)\)-\d+_\w+_\d+`,
	"5.MT_ARD":           `^.+?-(.+?)-\d+_\w+_\d+`,
	"11111.MT_ARD":       `^Filme_im_Ersten-(.+?)-\d+`,
	"6.MT_RBB":           `^.+?rbb-(.+?)-[0-9A-Fa-f\-]+_\d+`,
	"7.3_SAT":            `^(.+)_-_.+`,
	"998.Simple":         `^(.+)\s-\s.+\.\w\w\w$`,
	"999.Simple":         `^(.+)\.\w\w\w$`,
}

func (f *File) BeforeSave() (err error) {
	if f.FileName == "" {
		f.FileName = path.Base(f.FullPath)
	}
	return
}
func (f *File) Create(db *gorm.DB, tmdb TMDBClients) error {
	if err := db.Create(&f).Error; err != nil {
		log.Error(err)
		return err
	}
	movie := Movie{}
	//movie.FileID = f.ID
	regex := map[string]string{}
	//movie.Title = Translatename(f.FileName)
	movie.Title = Translatename(f.FileName, regex)
	//movie.WatchList = true
	//if s.Config.Mode != "testing" {
	err := movie.GetMeta(tmdb)
	if err != nil {
		log.Error(err)
	}
	//}
	//spew.Dump(movie)
	movie.FileID = f.ID

	if err := db.Create(&movie).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
func (f *File) Move(dir string) (string, error) {
	//target := viper.GetString("TargetDirectory")
	var newpath string
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		//	newpath = fmt.Sprintf("%s/%s/%s", target, dir, f.FileName)
		newpath = fmt.Sprintf("%s/%s", dir, f.FileName)
		err = os.Rename(f.FullPath, newpath)
		if err != nil {
			return "", fmt.Errorf("%s: %s", dir, err)
		}
		// path is a directory

	} else {
		return "", err
	}
	return newpath, nil
}

func Translatename(filename string, regexes map[string]string) string {
	if len(regexes) > 0 {
		log.Debug("merging regexes\n")
	}
	ret := ""
	keys := make([]string, 0)
	for k := range defaultRegexes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		//fmt.Println(k, defaultRegexes[k])
		name := k
		//for name, regex := range defaultRegexes {

		re := regexp.MustCompile(defaultRegexes[k])
		match := re.FindStringSubmatch(filename)
		if len(match) != 0 {
			i := len(match)
			match2 := match[1:i]
			ret = strings.Join(match2, " ")
			log.Debugf("Found %s Name: %s\n", name, ret)
			break
		}
	}
	ret = strings.Replace(ret, "_", " ", -1)
	ret = strings.TrimSpace(ret)
	return ret
}
