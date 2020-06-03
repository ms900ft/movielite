package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func ConnectDataBase() *gorm.DB {
	database, err := gorm.Open("sqlite3", "/tmp/test.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{}, &File{}, &Movie{}, &MovieSearchResults{}, &MovieShort{},
		&TMDBMovie{}, &Credits{}, &Cast{}, &Crew{}, &Genres{})

	return database
}
