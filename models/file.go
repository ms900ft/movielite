package models

import (
	"path"

	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type File struct {
	gorm.Model
	FullPath string `json:"fullPath" gorm:"Type:text"`
	FileName string `json:"fileName"`
}

func (f *File) BeforeSave() (err error) {
	if f.FileName == "" {
		f.FileName = path.Base(f.FullPath)
	}
	return
}
