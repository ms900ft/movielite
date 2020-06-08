package movielight

import (
	"fmt"
	"ms/movielight/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

type UpdateMovie struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	OrgName string `json:"org_name"`

	//Multiplechoice MovieSearchResults `json:"multiplechoice"`
	//File           File
	IsTv   bool `json:"is_tv"`
	Rating int  `json:"rating"`

	LastScanned time.Time
}
type Meta struct {
	Total int64 `json:"total"`
}

type Movielist struct {
	Data []models.Movie `json:"data"`
	Meta meta           `json:"meta"`
}

//Query movie query parameter
type Query struct {
	Orderby     string `form:"orderby"`
	Qtitel      string `form:"title"`
	Alpha       string `form:"alpha"`
	Genre       int64  `form:"genre"`
	Crew        int64  `form:"crew"`
	Person      int64  `form:"person"`
	Cast        int64  `form:"cast"`
	Country     string `form:"country"`
	LastScanned string `form:"last_scanned"`
	Limit       int64  `form:"limit,default=30"`
	Offset      int64  `form:"offset,default=0"`
	Show        string `form:"show"`
}

func (s *Service) getMovie(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var movie models.Movie

	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&movie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, movie)

}

func (s *Service) getMovies(c *gin.Context) {
	var q Query
	var fulltext = true
	err := c.Bind(&q)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	db := s.DB
	var movies []models.Movie
	var count int64

	tx := db.Set("gorm:auto_preload", true).Model(&models.Movie{}).
		Joins("JOIN files on files.id=movies.file_id").
		Where("is_tv = false")

	if len(q.Qtitel) > 0 {
		if fulltext {
			tx = tx.Joins("JOIN moviesearch on moviesearch.ID = movies.id").
				Where("moviesearch = ?", fmt.Sprintf("%s*", q.Qtitel))
			//	strings.Replace(q.Qtitel, " ", "&", -1)+":*")
		} else {
			tx = tx.Where("movies.title LIKE ?", fmt.Sprint("%", q.Qtitel, "%"))
		}
	}
	//xxx not working, test join with other parameter
	if q.LastScanned != "" {
		tx = tx.Where("last_scanned < ?", q.LastScanned).
			Joins("LEFT JOIN tmdb_movies on tmdb_movies.movie_id=movies.id")
		tx = tx.Where("tmdb_movies.id is not null")
	}

	switch {
	case q.Show == "multiple":
		tx = tx.Joins("JOIN movie_search_results on movie_search_results.movie_id=movies.id")
	case q.Show == "recent":
		tx = tx.Joins("LEFT JOIN recentlies on recentlies.movie_id = movies.id").
			Where("recentlies.user_id = ?", s.User.ID)
	case q.Show == "unrated":
		tx = tx.Where("rating = ?", 0)
	case q.Show == "notitle":
		tx = tx.Where("title = ?", "")
	case q.Show == "nodesc":
		tx = tx.Joins("LEFT JOIN tmdb_movies on tmdb_movies.movie_id=movies.id").
			Where("tmdb_movies.movie_id is null")

	}
	if q.Genre > 0 {
		tx = tx.Joins("JOIN tmdb_movies on tmdb_movies.movie_id=movies.id").
			Joins("JOIN genres ON genres.tmdb_movie_id = tmdb_movies.id").
			Where("genres.tmdb_id = ?", q.Genre)
	}
	if len(q.Country) > 0 {
		tx = tx.Joins("JOIN tmdb_movies on tmdb_movies.movie_id=movies.id").
			Joins("JOIN production_countries ON production_countries.tmdb_movie_id = tmdb_movies.id").
			Where("production_countries.iso3166_1 = ?", q.Country)
	}
	if q.Person > 0 {
		tx = tx.Joins("JOIN tmdb_movies on tmdb_movies.movie_id=movies.id").
			Joins("LEFT JOIN credits ON credits.id = tmdb_movies.id").
			Joins("LEFT JOIN casts ON casts.credits_id = credits.id").
			Joins("LEFT JOIN crews ON crews.credits_id = credits.id").
			Where("casts.cast_org_id = ? ", q.Person).
			Or("crews.crew_org_id = ?", q.Person)

	}

	//order
	switch {
	case q.Orderby == "name" || len(q.Alpha) > 0:
		tx = tx.Order("title ASC")
	case q.Orderby == "recent":
		tx = tx.Order("recentlies.last_played DESC")
	case q.Orderby == "last_scanned":
		tx = tx.Order("last_scanned")
	case q.Show == "watchlist":
		tx = tx.Order("watchlist to do")
	default:
		tx = tx.Group("movies.id")
		tx = tx.Order("files.created_at DESC")
	}

	tx = tx.Count(&count).
		Offset(q.Offset).
		Limit(q.Limit)

	if err := tx.Find(&movies).
		Error; err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ml := Movielist{}
	ml.Data = movies
	ml.Meta.Total = count
	c.JSON(http.StatusOK, ml)
}

func (s *Service) deleteMovie(c *gin.Context) {
	db := s.DB
	var movie models.Movie
	if err := db.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&movie)

	c.JSON(http.StatusOK, movie)
}

func (s *Service) createMovie(c *gin.Context) {
	db := s.DB
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f := []models.File{}
	if err := db.Where("full_path = ?", movie.File.FullPath).Find(&f).Error; err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(f) > 0 {
		content := "file already exists"
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
		return
	}
	if err := db.Create(&movie).Error; gorm.IsRecordNotFoundError(err) {
		content := gin.H{"error: ": "create" + err.Error()}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
		return
	}
	c.JSON(http.StatusCreated, movie)
}

func (s *Service) updateMovie(c *gin.Context) {
	db := s.DB
	var input UpdateMovie
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("movie binding input: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var movie models.Movie
	if err := db.Where("id = ?", input.ID).First(&movie).Error; err != nil {
		log.Errorf("movie not found: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	if err := db.Model(&movie).Update(input).Error; err != nil {
		log.Errorf("files update error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (s *Service) playMovie(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var movie models.Movie

	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&movie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	err := open.Start(movie.File.FullPath)
	if err != nil {
		log.Error(err)
	}

	recently := models.Recently{MovieID: movie.ID, UserID: s.User.ID, LastPlayed: time.Now()}
	if err := db.Save(&recently).Error; err != nil {
		content := gin.H{"error: ": "create" + err.Error()}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
