package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	//ID       uint `json:"id" gorm:"primary_key"`
	Username string
}
