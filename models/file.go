package models

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type File struct {
	gorm.Model
	FullPath string `json:"FullPath" gorm:"Type:text;UNIQUE;index"`
	FileName string `json:"FileName"`
	//CreatedAt time.Time
}

func (f *File) BeforeSave() (err error) {
	if f.FileName == "" {
		f.FileName = path.Base(f.FullPath)
	}
	return
}
func (f *File) Move(dir string) (string, error) {
	target := viper.GetString("TargetDirectory")
	var newpath string
	if stat, err := os.Stat(target); err == nil && stat.IsDir() {
		newpath = fmt.Sprintf("%s/%s/%s", target, dir, f.FileName)
		err = os.Rename(f.FullPath, newpath)
		if err != nil {
			return "", err
		}
		// path is a directory
	} else {
		return "", err
	}
	return newpath, nil
}
