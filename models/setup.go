package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func ConnectDataBase() *gorm.DB {
	database, err := gorm.Open("sqlite3", "/tmp/test.db")
	database.LogMode(true)
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{}, &File{}, &Movie{}, &MovieSearchResults{}, &MovieShort{},
		&TMDBMovie{}, &Credits{}, &Cast{}, &Crew{}, &Genres{}, &SpokenLanguages{},
		&ProductionCompanies{}, &ProductionCountries{}, &User{}, &Watchlist{}, &Recently{})
	_, err = database.DB().Exec("CREATE VIRTUAL  TABLE IF NOT EXISTS moviesearch USING fts5(ID, Title, Overview,Credits);")
	if err != nil {
		log.Fatal(err)
	}
	return database
}
