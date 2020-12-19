package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func ConnectDataBase(dbname string) *gorm.DB {
	database, err := gorm.Open("sqlite3", dbname)
	//database.LogMode(true)
	if err != nil {
		log.Panicf("Failed to connect to database! %s: %s", dbname, err)
	}

	database.AutoMigrate(&User{}, &File{}, &Movie{}, &MovieSearchResults{}, &MovieShort{},
		&TMDBMovie{}, &Credits{}, &Cast{}, &Crew{}, &Genres{}, &SpokenLanguages{},
		&ProductionCompanies{}, &ProductionCountries{}, &User{}, &Watchlist{}, &Recently{},
	)
	_, err = database.DB().Exec(`CREATE VIRTUAL  TABLE IF NOT EXISTS fulltexts
	USING fts5(movie_id, title, overview,credits);`)
	if err != nil {
		log.Fatal(err)
	}
	return database
}
