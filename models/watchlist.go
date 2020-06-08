package models

import (
	"time"
)

type Watchlist struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	MovieID   int64     `sql:"type:bigint REFERENCES movies(id) ON DELETE CASCADE"`
	UserID    int64     `sql:"type:bigint REFERENCES users(id) ON DELETE CASCADE"`
	CreatedAt time.Time `sql:"index"`
}
