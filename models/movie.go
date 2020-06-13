package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Movie struct {
	ID             int64     `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `sql:"index"`
	UpdatedAt      time.Time
	DeletedAt      *time.Time          `sql:"index"`
	FileID         uint                `json:"file_id"`
	Title          string              `json:"title"`
	OrgName        string              `json:"org_name"`
	Meta           *TMDBMovie          `json:"meta,omitempty"`
	Multiplechoice *MovieSearchResults `json:"multiplechoice,omitempty" gorm:"foreignkey:MovieID"`
	File           File
	IsTv           bool `json:"is_tv" gorm:"index"`
	Rating         int  `json:"rating"`
	Watchlist      bool `json:"watchlist" gorm:"-"`
	LastScanned    time.Time
}

type MovieSearchResults struct {
	gorm.Model
	//tmdb.MovieSearchResults
	MovieID      uint
	Page         int
	Results      []MovieShort `gorm:"foreignkey:MovieSearchResultsID"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

type MovieShort struct {
	gorm.Model
	MovieSearchResultsID uint
	Adult                bool   `json:"adult"`
	BackdropPath         string `json:"backdrop_path"`
	//ID                   int    `json:"id"`
	OriginalTitle string `json:"original_title"`
	//GenreIDs             []int32 `json:"genre_ids"`
	Popularity  float32 `json:"popularity"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	Video       bool    `json:"video"`
	VoteAverage float32 `json:"vote_average"`
	VoteCount   uint32  `json:"vote_count"`
}

type TMDBMovie struct {
	//ID           int       `gorm:"AUTO_INCREMENT"`
	CreatedAt    time.Time `gorm:"index"`
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
	TMDBMovieID  int        `json:"ID"`
	Adult        bool
	MovieID      uint   `gorm:"primary_key" json:"id"`
	BackdropPath string `json:"backdrop_path"`
	// BelongsToCollection bool   `json:"belongs_to_collection"`
	//BelongsToCollection CollectionShort `json:"belongs_to_collection"`
	Budget   uint32
	Genres   []Genres `json:"Genres" gorm:"many2many:movie_genres;ForeignKey:MovieId"`
	Homepage string
	//ID               int
	ImdbID              string `json:"imdb_id"`
	OriginalLanguage    string `json:"original_language"`
	OriginalTitle       string `json:"original_title"`
	Overview            string
	Popularity          float32
	PosterPath          string                `json:"poster_path"`
	ProductionCompanies []ProductionCompanies `json:"production_companies" gorm:"many2many:movie_production_companies;ForeignKey:MovieId"`
	ProductionCountries []ProductionCountries `json:"production_countries" gorm:"many2many:movie_production_countries;ForeignKey:MovieId"`
	ReleaseDate         string                `json:"release_date"`
	Revenue             uint32
	Runtime             uint32
	SpokenLanguages     []SpokenLanguages `json:"spoken_languages" gorm:"many2many:movie_spoken_languages;ForeignKey:MovieId"`
	Status              string
	Tagline             string
	Title               string
	Video               bool
	VoteAverage         float32 `json:"vote_average"`
	VoteCount           uint32  `json:"vote_count"`
	// AlternativeTitles *MovieAlternativeTitles `json:"alternative_titles,omitempty"`
	Credits Credits `json:",omitempty"`
	// Images            *MovieImages            `json:",omitempty"`
	// Keywords          *MovieKeywords          `json:",omitempty"`
	// Releases          *MovieReleases          `json:",omitempty"`
	// Videos            *MovieVideos            `json:",omitempty"`
	// Translations      *MovieTranslations      `json:",omitempty"`
	// Similar           *MoviePagedResults      `json:",omitempty"`
	// Reviews           *MovieReviews           `json:",omitempty"`
	// Lists             *MovieLists             `json:",omitempty"`
	// Changes           *MovieChanges           `json:",omitempty"`
	// Rating            *MovieRating            `json:",omitempty"`
	// ExternalIDs       *MovieExternalIds       `json:"external_ids,omitempty"`
}
type Credits struct {
	ID int64 `gorm:"primary_key"`

	TMDBMovieMovieID uint   `gorm:"index"`
	Crew             []Crew `gorm:"many2many:credits_crews;ForeignKey:TMDBMovieMovieID"`
	Cast             []Cast `gorm:"many2many:credits_casts;ForeignKey:TMDBMovieMovieID"`
}

type Cast struct {
	ID int64 `gorm:"primary_key"`
	//TMDBMovieMovieID uint  `gorm:"index"`
	//CreditsID        uint  `gorm:"index"`
	CastID    int `json:"cast_id"`
	CastOrgID int `json:"ID" gorm:"index"`

	CreditID    string `json:"credit_id"`
	Character   string
	Name        string
	XGender     int `json:"gender"`
	Order       int
	ProfilePath string `json:"profile_path"`
}

type Crew struct {
	ID int64 `gorm:"primary_key"`
	//TMDBMovieMovieID uint   `gorm:"index"`
	//CreditsID        uint   `gorm:"index"`
	CrewOrgID  int    `json:"ID" gorm:"index"`
	CreditID   string `json:"credit_id"`
	Department string
	Gender     int `json:"gender"`
	//	ID          int
	Job         string
	Name        string
	ProfilePath string `json:"profile_path"`
}

type Genres struct {
	ID      int64 `gorm:"primary_key"`
	Tmdb_id int   `json:"ID" gorm:"index:genreid"`
	//TMDBMovieID uint
	Name string
}

type SpokenLanguages struct {
	Iso639_1 string `gorm:"primary_key" json:"iso_639_1"`
	Name     string
}

type ProductionCompanies struct {
	ID        int64 `gorm:"primary_key" json:"id"`
	Name      string
	LogoPath  string `json:"logo_path"`
	Iso3166_1 string `json:"origin_country"`
}

type ProductionCountries struct {

	//TMDBMovieMovieID uint
	Iso3166_1 string `gorm:"primary_key" json:"iso_3166_1"`
	Name      string
}

const (
	//DEFLANG   = "de-DE"
	WATCHLIST = "watchlist"
)

func (m *Movie) GetMeta() (err error) {
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("can't get metadata")
		}
	}()
	lang := viper.GetString("language")
	apikey := viper.GetString("TMDB.ApiKey")
	conf := tmdb.Config{APIKey: apikey}
	TMDb := tmdb.Init(conf)
	var options = make(map[string]string)
	options["language"] = lang
	res, err := TMDb.SearchMovie(m.Title, options)
	if err != nil {
		log.Error(err)
	}
	//spew.Dump(res)
	hit := 0
	hits := 0
	for _, search := range res.Results {
		if search.Title == m.Title {
			hit = search.ID
			hits++
		}
	}
	if len(res.Results) == 1 {
		hit = res.Results[0].ID
	}
	if hit == 0 || hits > 1 {
		log.Debug("Multiple results")

		//jsonRes := []byte("{}") //Default value in case of error
		//jsonRes, _ := json.MarshalIndent(res, "", "  ")
		//var mc MovieSearchResults
		//err = json.Unmarshal(jsonRes, &mj)
		//spew.Dump(m)
		//m.Meta = nil
		js, err := json.Marshal(res)
		if err != nil {
			log.Error(err)
			return err
		}
		msr := MovieSearchResults{}
		err = json.Unmarshal(js, &msr)
		if err != nil {
			log.Error(err)
			return err
		}
		if res.TotalResults > 0 {
			m.Multiplechoice = &msr
		}
	} else {
		meta, err := getTMDBMeta(hit)
		m.Meta = &meta
		//m.Multiplechoice = nil
		if err != nil {
			log.Error(err)
		}
	}
	return err
}

func getTMDBMeta(id int) (TMDBMovie, error) {
	apikey := viper.GetString("TMDB.ApiKey")
	lang := viper.GetString("language")
	conf := tmdb.Config{APIKey: apikey}
	TMDb := tmdb.Init(conf)
	var options = make(map[string]string)
	options["append_to_response"] = "credits"
	options["language"] = lang
	res, err := TMDb.GetMovieInfo(id, options)
	if err != nil {
		log.Error(err)
		return TMDBMovie{}, err
	}
	//spew.Dump(res)
	// err = preFetchImages(res)
	// if err != nil {
	// 	log.Error(err)
	// }
	//spew.Dump(res)
	//jsonRes := []byte("{}") //Default value in case of error
	js, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		return TMDBMovie{}, err
	}
	msr := TMDBMovie{}
	err = json.Unmarshal(js, &msr)
	if err != nil {
		log.Error(err)
		return msr, err
	}
	return msr, err
}

func (m *Movie) MetaById(metaid int) error {
	meta, err := getTMDBMeta(metaid)
	if err != nil {
		log.Error(err)
		return err
	}
	m.Meta = &meta
	m.Multiplechoice = nil
	return nil
}

func (m *Movie) AfterCreate(scope *gorm.Scope) (err error) {

	err = scope.DB().Exec("INSERT INTO moviesearch (ID,Title,Overview,Credits) VALUES($1, $2, $3, $4)",
		m.ID, m.Title, m.Meta.Overview, m.GetCredits()).Error
	if err != nil {
		log.Error(err)
		return err
	}

	var users []User
	if err := scope.DB().Find(&users).Error; err != nil {
		log.Error(err)
		return err
	}
	for _, user := range users {
		w := Watchlist{UserID: user.ID, MovieID: m.ID}
		if err := scope.DB().Create(&w).Error; gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return err
		}
	}
	return
}

func (m *Movie) AfterUpdate(scope *gorm.Scope) (err error) {
	db := scope.DB()
	if err := db.Model(&m.Meta).Association("ProductionCountries").
		Replace(m.Meta.ProductionCountries).Error; err != nil {
		log.Errorf("ProductionCountries update error: %s", err)

		return err
	}
	if err := db.Model(&m.Meta).Association("ProductionCompanies").
		Replace(m.Meta.ProductionCompanies).Error; err != nil {
		log.Errorf("ProductionCompanies update error: %s", err)
		return err
	}
	if err := db.Model(&m.Meta).Association("Genres").
		Replace(m.Meta.Genres).Error; err != nil {
		log.Errorf("Genres update error: %s", err)
		return err
	}
	if err := db.Model(&m.Meta).Association("SpokenLanguages").
		Replace(m.Meta.SpokenLanguages).Error; err != nil {
		log.Errorf("SpokenLanguages update error: %s", err)
		return err
	}
	if err := db.Model(&m.Meta.Credits).Association("Crew").
		Replace(m.Meta.Credits.Crew).Error; err != nil {
		log.Errorf("Crew update error: %s", err)
		return err
	}
	if err := db.Model(&m.Meta.Credits).Association("Cast").
		Replace(m.Meta.Credits.Cast).Error; err != nil {
		log.Errorf("Cast update error: %s", err)

		return err
	}

	err = db.Exec("UPDATE moviesearch set Title =$2,Overview =$3,Credits=$4 WHERE ID = $1",
		m.ID, m.Title, m.Meta.Overview, m.GetCredits()).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return
}

func (m *Movie) GetCredits() string {
	var credits string
	if m.Meta != nil {
		for _, c := range m.Meta.Credits.Cast {
			credits += c.Name + " "
		}
		for _, c := range m.Meta.Credits.Crew {
			credits += c.Name + " "
		}
	}
	return credits
}
