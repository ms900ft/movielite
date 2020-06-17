package models

type Fulltext struct {
	MovieID  int64 `gorm:"primary_key"`
	Title    string
	Overview string
	Credits  string
}
