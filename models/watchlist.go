package models

import (
	"time"
)

type Watchlist struct {
	MovieID   int64 `gorm:"primary_key" sql:"type:bigint REFERENCES movies(id) ON DELETE CASCADE"`
	UserID    int64 `gorm:"primary_key" sql:"type:bigint REFERENCES users(id) ON DELETE CASCADE"`
	CreatedAt time.Time
}
