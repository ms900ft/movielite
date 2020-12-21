package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bios-Marcel/wastebasket"
	"github.com/gammazero/workerpool"
	"github.com/jinzhu/gorm"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Movie struct {
	ID                   int64     `json:"id" gorm:"primary_key"`
	CreatedAt            time.Time `sql:"index"`
	UpdatedAt            time.Time
	DeletedAt            *time.Time `sql:"index"`
	FileID               uint       `json:"file_id"`
	TMDBMovieID          uint       `sql:"index"`
	MovieSearchResultsID uint
	Title                string              `json:"title" binding:"required"`
	OrgName              string              `json:"org_name"`
	Meta                 *TMDBMovie          `json:"meta" gorm:"foreignkey:TMDBMovieID"`
	Multiplechoice       *MovieSearchResults `json:"multiplechoice" gorm:"foreignkey:MovieSearchResultsID"`
	File                 File                `json:"File" binding:"required"`
	IsTv                 bool                `json:"is_tv" gorm:"index"`
	Rating               int                 `json:"rating"`
	Watchlist            bool                `json:"watchlist" gorm:"-"`
	LastScanned          time.Time
	Deleted              bool `gorm:"-"`
}

type MovieSearchResults struct {
	ID        int       `gorm:"AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	Page         int
	Results      []MovieShort `json:"Results" gorm:"many2many:movie_search_results_movie_short"`
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
	ID        int       `gorm:"AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	//TMDBMovieID int        `json:"ID"`
	Adult bool
	//MovieID      int64  `gorm:"primary_key" json:"id"`
	BackdropPath string `json:"backdrop_path"`
	// BelongsToCollection bool   `json:"belongs_to_collection"`
	//BelongsToCollection CollectionShort `json:"belongs_to_collection"`
	Budget   uint32
	Genres   []Genres `json:"Genres" gorm:"many2many:tmdb_movie_genres"`
	Homepage string
	//ID               int
	ImdbID              string `json:"imdb_id"`
	OriginalLanguage    string `json:"original_language"`
	OriginalTitle       string `json:"original_title"`
	Overview            string
	Popularity          float32
	PosterPath          string                `json:"poster_path"`
	ProductionCompanies []ProductionCompanies `json:"production_companies" gorm:"many2many:tmdb_movie_production_companies"`
	ProductionCountries []ProductionCountries `json:"production_countries" gorm:"many2many:tmdb_movie_production_countries"`
	ReleaseDate         string                `json:"release_date"`
	Revenue             uint32
	Runtime             uint32
	SpokenLanguages     []SpokenLanguages `json:"spoken_languages" gorm:"many2many:tmdb_movie_spoken_languages"`
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

	TMDBMovieID uint   `gorm:"index"`
	Crew        []Crew `gorm:"many2many:credits_crews"`
	Cast        []Cast `gorm:"many2many:credits_casts"`
}

type Cast struct {
	//ID int64 `gorm:"primary_key"`
	//TMDBMovieMovieID uint  `gorm:"index"`
	//CreditsID        uint  `gorm:"index"`
	//CastID    int `json:"cast_id"`
	PersonID int `json:"ID" gorm:"index"`

	ID          string `json:"credit_id" gorm:"index"`
	Character   string
	Name        string
	XGender     int `json:"gender"`
	Order       int
	ProfilePath string `json:"profile_path"`
}

type Crew struct {
	//ID int64 `gorm:"primary_key"`
	//TMDBMovieMovieID uint   `gorm:"index"`
	//CreditsID        uint   `gorm:"index"`
	PersonID   int    `json:"ID" gorm:"index"`
	ID         string `json:"credit_id" gorm:"index"`
	Department string
	Gender     int `json:"gender"`
	//	ID          int
	Job         string
	Name        string
	ProfilePath string `json:"profile_path"`
}

type Genres struct {
	//ID      int64 `gorm:"primary_key"`
	TmdbID int `gorm:"primary_key" json:"ID" `
	//
	TMDBMovieID uint `gorm:"index"`
	Name        string
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

const raw = `
on run argv
  tell application "Finder"
    repeat with f in argv
      move (f as POSIX file) to trash
    end repeat
  end tell
end run
`

func (m *Movie) GetMeta(t TMDBClients) (err error) {
	// defer func() {
	// 	// recover from panic if one occurred. Set err to nil otherwise.
	// 	if recover() != nil {
	// 		err = errors.New("can't get metadata for " + m.Title)
	// 	}
	//	}()
	lang := viper.GetString("language")
	//apikey := viper.GetString("TMDB.ApiKey")
	//conf := tmdb.Config{APIKey: apikey}
	//TMDb := tmdb.Init(conf)
	var options = make(map[string]string)
	options["language"] = lang
	res, err := t.SearchMovie(m.Title, options)
	if err != nil {
		log.Error(err)
	}
	if res == nil {
		m.Multiplechoice = nil
		m.Meta = nil
		return nil
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
			m.Meta = nil
		}
	} else {
		meta, err := getTMDBMeta(t, nil, hit)
		m.Meta = &meta
		m.Multiplechoice = nil
		//m.Multiplechoice = nil
		if err != nil {
			log.Error(err)
		}
	}
	return err
}

func getTMDBMeta(t TMDBClients, wp *workerpool.WorkerPool, id int) (TMDBMovie, error) {
	//apikey := viper.GetString("TMDB.ApiKey")
	lang := viper.GetString("language")
	//conf := tmdb.Config{APIKey: apikey}
	//TMDb := tmdb.Init(conf)
	var options = make(map[string]string)
	options["append_to_response"] = "credits"
	options["language"] = lang
	res, err := t.GetMovieInfo(id, options)
	if err != nil {
		log.Error(err)
		return TMDBMovie{}, err
	}
	//spew.Dump(res)
	err = preFetchImages(res, wp)
	if err != nil {
		log.Error(err)
		return TMDBMovie{}, err
	}
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

func (m *Movie) MetaByID(t TMDBClients, wp *workerpool.WorkerPool, metaid int) error {
	meta, err := getTMDBMeta(t, wp, metaid)
	if err != nil {
		log.Error(err)
		return err
	}
	m.Meta = &meta
	m.Multiplechoice = nil
	//m.Meta.MovieID = m.ID
	return nil
}

func (m *Movie) AfterCreate(scope *gorm.Scope) (err error) {
	if m.DeletedAt != nil {
		return
	}
	fulltext := Fulltext{MovieID: m.ID}
	if m.Meta != nil {
		if m.Title != m.Meta.OriginalTitle {
			fulltext.Title = fulltext.Title + " " + m.Meta.OriginalTitle
		}
	}
	if m.Meta != nil {
		fulltext.Overview = m.Meta.Overview
		fulltext.Credits = m.GetCredits()
	}
	if err := scope.DB().Create(&fulltext).Error; err != nil {
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

func (m *Movie) FullTextIndex(tx *gorm.DB) error {
	fulltext := Fulltext{MovieID: m.ID}
	found := true
	if err := tx.Where("movie_id = ?", m.ID).First(&fulltext).Error; gorm.IsRecordNotFoundError(err) {
		found = false
	}
	fulltext.Title = m.Title
	if m.Meta != nil {
		if m.Title != m.Meta.OriginalTitle {
			fulltext.Title = fulltext.Title + " " + m.Meta.OriginalTitle
		}
		fulltext.Overview = m.Meta.Overview
		fulltext.Credits = m.GetCredits()
	}
	if found {
		if err := tx.Model(&fulltext).Update(&fulltext).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Create(&fulltext).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *Movie) AfterUpdate(scope *gorm.Scope) (err error) {
	return m.FullTextIndex(scope.DB())

}

func (m *Movie) AfterDelete(scope *gorm.Scope) (err error) {
	var users []User
	if err := scope.DB().Find(&users).Error; err != nil {
		log.Error(err)
		return err
	}
	for _, user := range users {
		w := Watchlist{UserID: user.ID, MovieID: m.ID}
		if err := scope.DB().Delete(&w).Error; err != nil {
			log.Error(err)
			return err
		}
		r := Recently{UserID: user.ID, MovieID: m.ID}
		if err := scope.DB().Delete(&r).Error; err != nil {
			log.Error(err)
			return err
		}
	}
	f := Fulltext{MovieID: m.ID}
	if err := scope.DB().Delete(&f).Error; err != nil {
		log.Error(err)
		return err
	}
	_, e := os.Stat(m.File.FullPath)
	if !os.IsNotExist(e) {
		trashcan := viper.GetBool("TrashCan")
		if m.File.FullPath != "" {
			if trashcan {
				log.Debugf("moving %s to trashcan", m.File.FullPath)
				err := wastebasket.Trash(m.File.FullPath)
				if err != nil {
					log.Error(err)
					return err
				}
			} else {
				log.Debugf("deleting movie %s", m.File.FullPath)
				err := os.Remove(m.File.FullPath)
				if err != nil {
					log.Errorf("cannot delete file %s: %s", m.File.FullPath)
					return err
				}
			}
		}
	}
	return err
}

func (m *Movie) GetCredits() string {
	var credits string
	//if m.Meta != nil {
	for _, c := range m.Meta.Credits.Cast {
		credits += c.Name + " "
	}
	for _, c := range m.Meta.Credits.Crew {
		credits += c.Name + " "
	}
	//	}
	return credits
}

func (m *Movie) DeleteMeta(db *gorm.DB, movie *Movie) (err error) {
	if err := db.Model(movie.Meta).Association("ProductionCountries").
		Delete(movie.Meta.ProductionCountries).Error; err != nil {
		log.Errorf("ProductionCountries delete error: %s", err)
		return err
	}
	if err := db.Model(movie.Meta).Association("ProductionCompanies").
		Delete(movie.Meta.ProductionCompanies).Error; err != nil {
		log.Errorf("ProductionCompanies delete error: %s", err)
		return err
	}
	if err := db.Model(movie.Meta).Association("Genres").
		Delete(movie.Meta.Genres).Error; err != nil {
		log.Errorf("Genres update error: %s", err)
		return err
	}
	if err := db.Model(movie.Meta).Association("SpokenLanguages").
		Delete(movie.Meta.SpokenLanguages).Error; err != nil {
		log.Errorf("SpokenLanguages update error: %s", err)
		return err
	}

	if err := db.Model(movie.Meta.Credits).Association("Crew").
		Delete(movie.Meta.Credits.Crew).Error; err != nil {
		log.Errorf("Crew update error: %s", err)
		return err
	}
	if err := db.Model(movie.Meta.Credits).Association("Cast").
		Delete(movie.Meta.Credits.Cast).Error; err != nil {
		log.Errorf("Cast update error: %s", err)

		return err
	}
	if err := db.Model(movie.Meta).Association("Credits").
		Delete(movie.Meta.Credits).Error; err != nil {
		log.Errorf("Crew update error: %s", err)
		return err
	}
	return err
}

func Trash(f string, trash string) (trashcan string, err error) {
	bin, err := exec.LookPath("osascript")
	if err != nil {
		err = fmt.Errorf("not yet supported")
		return
	}

	if _, err = os.Stat(trash); err != nil {
		err = fmt.Errorf("trash %s not found", trash)
		return "", err
	}

	path, err := filepath.Abs(f)
	if err != nil {
		return
	}
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	_ = name

	dest := filepath.Join(trash, base)
	if _, err = os.Stat(dest); err == nil {
		err = fmt.Errorf("already exists")
		return
	}
	trashcan = dest
	log.Debug(path)
	params := append([]string{"-e", raw}, path)
	cmd := exec.Command(bin, params...)
	log.Debugf("%+v", params)
	if err = cmd.Run(); err != nil {
		log.Error(err)
		return
	}

	if _, err = os.Stat(trashcan); err != nil {
		trashcan = ""
	}

	return trashcan, err
}

func preFetchImages(movie *tmdb.Movie, wp *workerpool.WorkerPool) error {
	urls := preFetchURLS(movie)
	for _, url := range urls {
		log.Debugf("getting images %s", url)
		url := url
		if wp != nil {
			wp.Submit(func() {
				log.Debugf("Add url to queue: %s", url)
				time.Sleep(100 * time.Millisecond)

				client := &http.Client{}
				req, err := http.NewRequest("GET", url, nil)
				//req.Close = true
				if err != nil {
					log.Error(err)
					//return err
				}
				response, err := client.Do(req)

				if err != nil {
					log.Error(err)
					//return err
				}
				if response.StatusCode != http.StatusOK {
					log.Errorf(fmt.Sprintf("Can't get %s: %s Status %d", url, err, response.StatusCode))
					//return err
				}
				defer response.Body.Close()
			})
		}
	}
	return nil
}

func preFetchURLS(movie *tmdb.Movie) []string {
	port := viper.GetInt("Port")

	urls := []string{}
	baseurl := fmt.Sprintf("http://localhost:%d", port)
	//posterpath
	if movie.PosterPath != "" {
		sizes := []string{"w342", "w780", "w185", "w92"}
		for _, size := range sizes {
			url := fmt.Sprintf("%s/images/%s%s", baseurl, size, movie.PosterPath)
			urls = append(urls, url)
		}
	}
	if movie.BackdropPath != "" {
		url := fmt.Sprintf("%s/images/w300%s", baseurl, movie.BackdropPath)
		urls = append(urls, url)
	}
	if movie.Credits.Cast != nil {
		for _, cast := range movie.Credits.Cast {
			if cast.ProfilePath != "" {
				sizes := []string{"w45", "w500"}
				for _, size := range sizes {
					url := fmt.Sprintf("%s/images/%s%s", baseurl, size, cast.ProfilePath)
					urls = append(urls, url)
				}
			}
		}
	}
	if movie.Credits.Crew != nil {
		for _, crew := range movie.Credits.Crew {
			if crew.ProfilePath != "" {
				sizes := []string{"w45", "w500"}
				for _, size := range sizes {
					url := fmt.Sprintf("%s/images/%s%s", baseurl, size, crew.ProfilePath)
					urls = append(urls, url)
				}
			}
		}
	}
	//spew.Dump(urls)
	return urls
}
