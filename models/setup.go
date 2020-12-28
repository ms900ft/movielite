package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const adminName = "admin"

type DBConfig struct {
	DBName               string
	InitialAdminPassword string
}

func ConnectDataBase(c DBConfig) *gorm.DB {
	database, err := gorm.Open("sqlite3", c.DBName)
	//database.LogMode(true)
	if err != nil {
		log.Panicf("Failed to connect to database! %s: %s", c.DBName, err)
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
	err = addAdmin(database, c.InitialAdminPassword)
	if err != nil {
		log.Error(err)
	}

	return database
}

func addAdmin(db *gorm.DB, pass string) error {
	var user User
	rows, err := db.Where("user_name = ?", adminName).First(&user).Rows()
	if err != nil {
		return err
	}
	if !rows.Next() {
		p, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u := User{UserName: adminName, Password: string(p)}
		if err := db.Create(&u).Error; gorm.IsRecordNotFoundError(err) {
			return err
		}
	}
	log.Debug("admin created")
	return nil
}
