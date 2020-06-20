package models

type User struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	UserName string `gorm:"UNIQUE;index"`
}
