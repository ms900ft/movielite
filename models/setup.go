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

	// Performance indexes for common queries
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_movies_tmdb_movie_id ON movies(tmdb_movie_id)`,
		`CREATE INDEX IF NOT EXISTS idx_movies_file_id ON movies(file_id)`,
		`CREATE INDEX IF NOT EXISTS idx_movies_search_results_id ON movies(movie_search_results_id)`,
		`CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title)`,
		`CREATE INDEX IF NOT EXISTS idx_watchlists_user_movie ON watchlists(user_id, movie_id)`,
		`CREATE INDEX IF NOT EXISTS idx_recentlies_user_movie ON recentlies(user_id, movie_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tmdb_movie_genres_tmdb_id ON tmdb_movie_genres(genres_tmdb_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tmdb_movie_genres_movie_id ON tmdb_movie_genres(tmdb_movie_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tmdb_movie_countries_iso ON tmdb_movie_production_countries(production_countries_iso3166_1)`,
		`CREATE INDEX IF NOT EXISTS idx_tmdb_movie_countries_movie_id ON tmdb_movie_production_countries(tmdb_movie_id)`,
		`CREATE INDEX IF NOT EXISTS idx_credits_casts_cast_id ON credits_casts(cast_id)`,
		`CREATE INDEX IF NOT EXISTS idx_credits_casts_credits_id ON credits_casts(credits_id)`,
		`CREATE INDEX IF NOT EXISTS idx_credits_crews_crew_id ON credits_crews(crew_id)`,
		`CREATE INDEX IF NOT EXISTS idx_credits_crews_credits_id ON credits_crews(credits_id)`,
		`CREATE INDEX IF NOT EXISTS idx_casts_person_id ON casts(person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_crews_person_id ON crews(person_id)`,
		`CREATE INDEX IF NOT EXISTS idx_credits_tmdb_movie_id ON credits(tmdb_movie_id)`,
	}
	for _, idx := range indexes {
		if _, err := database.DB().Exec(idx); err != nil {
			log.Warnf("Index creation failed: %s", err)
		}
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
		u := User{UserName: adminName, Password: string(p), IsAdmin: true}
		if err := db.Create(&u).Error; gorm.IsRecordNotFoundError(err) {
			return err
		}
	}
	log.Debug("admin created")
	return nil
}
