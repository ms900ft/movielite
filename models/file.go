package models

import (
	"path"

	"github.com/jinzhu/gorm"
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
