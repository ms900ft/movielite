package models

import (
	"time"
)

type Recently struct {
	//ID        int64     `json:"id" gorm:"primary_key"`
	MovieID    int64 `gorm:"primary_key" sql:"type:bigint REFERENCES movies(id) ON DELETE CASCADE"`
	UserID     int64 `gorm:"primary_key" sql:"type:bigint REFERENCES users(id) ON DELETE CASCADE"`
	LastPlayed time.Time
}
