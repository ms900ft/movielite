package movielite

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ms900ft/movielite/models"
	"github.com/ms900ft/movielite/pkg/filemanager"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
)

type UpdateMovie struct {
	ID      int64  `json:"id"`
	Title   string `json:"title" binding:"required"`
	OrgName string `json:"org_name"`

	//Multiplechoice MovieSearchResults `json:"multiplechoice"`
	//File           File
	IsTv        bool              `json:"is_tv"`
	Rating      int               `json:"rating"`
	Watchlist   bool              `json:"watchlist"`
	Meta        *models.TMDBMovie `json:"meta"`
	LastScanned time.Time
}
type Meta struct {
	Total int64 `json:"total"`
}

type Movielist struct {
	Data []models.Movie `json:"data"`
	Meta meta           `json:"meta"`
}

// Query movie query parameter
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

// getMovie godoc
// @Summary Get a single movie
// @Description get a movie by ID
// @Tags movies
// @Produce  json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.Movie
// @Router /api/movie/{id} [get]
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

// getMovies godoc
// @Summary Get a list of movies
// @Description get movies
// @Tags movies
// @Produce  json
// @Param orderby query string false "Order by (name, recent, last_scanned)"
// @Param title query string false "Title filter"
// @Param alpha query string false "Alpha filter"
// @Param genre query int false "Genre ID filter"
// @Param crew query int false "Crew person ID filter"
// @Param person query int false "Person (cast/crew) ID filter"
// @Param cast query int false "Cast person ID filter"
// @Param country query string false "Country code filter"
// @Param show query string false "Show (multiple, unrated, duplicate, notitle, nodesc, watchlist)"
// @Success 200 {object} Movielist
// @Router /api/movie [get]
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
		Select(`movies.id, movies.file_id,movies.tmdb_movie_id,movies.movie_search_results_id,
		movies.title,movies.is_tv,movies.rating,
		CASE WHEN watchlists.movie_id is not null then true else false  end as watchlist`)
	if len(q.Qtitel) > 0 && fulltext {
		//	tx = tx.Joins("JOIN fulltexts on fulltexts.movie_id = movies.id").
		//		Where("fulltexts = ?", fmt.Sprintf("%s*", q.Qtitel))
		//	tx = tx.Where("movies.id IN (SELECT movie_id FROM fulltexts WHERE fulltexts MATCH ? order by bm25(fulltexts, 10.0, 5.0) )",
		//
		//		fmt.Sprintf("%s*", q.Qtitel))
		tx = tx.Joins("JOIN fulltexts on fulltexts.movie_id =movies.id ")
		tx = tx.Where("fulltexts match ?", fmt.Sprintf("%s*", q.Qtitel))
		//	strings.Replace(q.Qtitel, " ", "&", -1)+":*")
	}

	tx = tx.Joins("JOIN files on files.id=movies.file_id").
		Joins("Left Join watchlists ON (movies.id = watchlists.movie_id AND watchlists.user_id = ?)", s.Token.UserID).
		//Where("watchlists.user_id = ?", s.User.ID).
		Where("is_tv = false")

	if len(q.Qtitel) > 0 && !fulltext {
		tx = tx.Where("movies.title LIKE ?", fmt.Sprint("%", q.Qtitel, "%"))
	}

	//xxx not working, test join with other parameter
	if q.LastScanned != "" {
		tx = tx.Where("last_scanned < ?", q.LastScanned).
			Joins("LEFT JOIN tmdb_movies on tmdb_movies.movie_id=movies.id")
		tx = tx.Where("tmdb_movies.id is not null")
	}

	switch {
	case q.Show == "multiple":
		tx = tx.Where("movies.movie_search_results_id > 0")
	case q.Orderby == "recent":
		tx = tx.Joins("LEFT JOIN recentlies on recentlies.movie_id = movies.id").
			Where("recentlies.user_id = ?", s.Token.UserID)
	case q.Show == "unrated":
		tx = tx.Where("rating = ?", 0)
	case q.Show == "duplicate":
		tx = tx.Where("(select  count(*)  from movies inr  where inr.title = movies.title)> 1")
	case q.Show == "notitle":
		tx = tx.Where("title = ?", "")
	case q.Show == "nodesc":
		tx = tx.Where("movies.tmdb_movie_id == 0")
	case q.Show == "watchlist":
		tx = tx.Where("watchlists.movie_id is not null AND watchlists.user_id = ?", s.Token.UserID)
	}
	if q.Genre > 0 {
		tx = tx.Where(`tmdb_movie_id in (SELECT id
                    FROM tmdb_movies where (
										id in (select tmdb_movie_id
                    from tmdb_movie_genres
										where genres_tmdb_id
                    = ?))) `, q.Genre)
	}
	if len(q.Country) > 0 {
		tx = tx.Where(`tmdb_movie_id in (SELECT id
                    FROM tmdb_movies where (
										id in (select tmdb_movie_id
                    from tmdb_movie_production_countries
										where production_countries_iso3166_1
                    = ?))) `, q.Country)
	}
	if q.Person > 0 {
		tx = tx.Where(` tmdb_movie_id
     								IN (SELECT id
                    FROM tmdb_movies where (
                    id in (select tmdb_movie_id from credits WHERE
										id in (select credits_id from credits_casts WHERE
										cast_id in (select id from casts where
										person_id = ?)))
										OR
										id in (select tmdb_movie_id from credits WHERE
										id in (select credits_id from credits_crews WHERE
										crew_id in (select id from crews where
										person_id = ?)))
                    ))`, q.Person, q.Person)
	}
	tx = tx.Where("movies.is_tv = false")
	//order
	switch {
	case fulltext && len(q.Qtitel) > 0:
		tx = tx.Order("bm25(fulltexts, 1.0, 50.0, 5.0, 10.0)")

	case q.Orderby == "name" || len(q.Qtitel) > 0:
		tx = tx.Order("movies.title ASC")

	case q.Orderby == "recent":
		tx = tx.Order("recentlies.last_played DESC")
	case q.Orderby == "last_scanned":
		tx = tx.Order("movies.last_scanned")
	case q.Show == "watchlist":
		tx = tx.Order("watchlists.created_at DESC")
	case q.Show == "duplicate":
		tx = tx.Order("movies.title DESC")
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

// deleteMovie godoc
// @Summary Delete a movie
// @Description delete a movie by ID
// @Tags movies
// @Produce  json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.Movie
// @Router /api/movie/{id} [delete]
func (s *Service) deleteMovie(c *gin.Context) {
	db := s.DB
	var movie models.Movie
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&movie)

	c.JSON(http.StatusOK, movie)
}

// createMovie godoc
// @Summary Create a movie
// @Description create a new movie entry
// @Tags movies
// @Accept  json
// @Produce  json
// @Param movie body models.Movie true "Movie object"
// @Success 201 {object} models.Movie
// @Router /api/movie [post]
func (s *Service) createMovie(c *gin.Context) {
	db := s.DB
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if movie.Deleted {
		t := time.Now()
		movie.DeletedAt = &t
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

// updateMovie godoc
// @Summary Update a movie
// @Description update a movie by ID
// @Tags movies
// @Accept  json
// @Produce  json
// @Param id path int true "Movie ID"
// @Param movie body models.Movie true "Movie object"
// @Success 200 {object} models.Movie
// @Router /api/movie/{id} [put]
func (s *Service) updateMovie(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var input models.Movie
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("movie binding input: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var movie models.Movie
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&movie).
		Error; err != nil {
		log.Errorf("movie not found: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	old := movie
	movie.IsTv = input.IsTv
	if input.Title != old.Title {
		movie.Title = input.Title
		if err := movie.GetMeta(s.TMDBClient); err != nil {
			log.Errorf("movie meta update error: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		}
		if movie.Meta != nil {
			movie.Title = movie.Meta.Title
			movie.Multiplechoice = nil
		} else {
			movie.Meta = nil
		}
	}

	if err := db.Model(&old).Update(movie).Error; err != nil {
		log.Errorf("files update error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}
	if movie.Multiplechoice == nil {
		if err := db.Model(&old).Association("Multiplechoice").
			Clear().Error; err != nil {
			log.Errorf("remove multiplechoice: %s", err)
		}
	}
	if movie.Meta == nil {
		if err := db.Model(&old).Association("Meta").
			Clear().Error; err != nil {
			log.Errorf("remove meta: %s", err)
		}
	}
	hasWatchlist := true
	w := models.Watchlist{UserID: s.Token.UserID, MovieID: movie.ID}
	if db.Find(&w).First(&w).RecordNotFound() {
		hasWatchlist = false
	}

	if input.Watchlist != hasWatchlist {
		if hasWatchlist {
			if err := db.Delete(&w).Error; err != nil {
				log.Errorf("toggle watchlist error: %s", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
			}
			movie.Watchlist = false
		} else {
			if err := db.Create(&w).Error; err != nil {
				log.Errorf("toggle watchlist error: %s", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
			}
			movie.Watchlist = true
		}
	}
	c.JSON(http.StatusOK, movie)
}

// playMovie godoc
// @Summary Play a movie
// @Description play a movie file on the server
// @Tags movies
// @Produce  json
// @Param id path int true "Movie ID"
// @Param showdir query string false "Show directory in file manager instead of playing"
// @Router /api/movie/{id}/play [put]
func (s *Service) playMovie(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var movie models.Movie

	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&movie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	dir := c.Query("showdir")
	//works only on  macos
	if dir != "" {
		err := filemanager.ShowDir(movie.File.FullPath)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		c.JSON(http.StatusNoContent, nil)
		return
	}

	err := open.RunWith(movie.File.FullPath, s.Config.Player)
	if err != nil {
		log.Error(err)
	}

	recently := models.Recently{MovieID: movie.ID, UserID: s.Token.UserID, LastPlayed: time.Now()}
	if err := db.Save(&recently).Error; err != nil {
		content := gin.H{"error: ": "create" + err.Error()}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// addMeta godoc
// @Summary Add/Update metadata for a movie
// @Description Associates TMDB metadata with a movie
// @Tags movies
// @Produce  json
// @Param id path int true "Movie ID"
// @Param metaid path int true "TMDB Movie ID"
// @Success 200 {object} models.Movie
// @Router /api/movie/{id}/addMeta/{metaid} [put]
func (s *Service) addMeta(c *gin.Context) {
	db := s.DB
	id := c.Param("id")
	metaID, err := strconv.Atoi(c.Param("metaid"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var movie models.Movie

	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&movie).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	old := movie
	if metaID == 0 {
		movie.Meta = nil
		movie.TMDBMovieID = 0
		if err := db.Model(&old).Update(movie).Error; err != nil {
			log.Errorf("Movie update error: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
			return
		}
		if err := db.Model(&old).Association("Meta").
			Clear().Error; err != nil {
			log.Errorf("remove meta error: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
			return
		}
		c.JSON(http.StatusOK, movie)
		return
	}

	err = movie.MetaByID(s.TMDBClient, s.WorkerPool, metaID)
	movie.Title = movie.Meta.Title

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	//	if err := db.Debug().Model(&old).Association("TMDBMovie").Append(movie.Meta).Error; err != nil {
	if err := db.Model(&old).Update(movie).Error; err != nil {
		log.Errorf("Movie update error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}
	if err := db.Model(&old).Association("Multiplechoice").
		Clear().Error; err != nil {
		log.Errorf("remove multiplechoice error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}

	c.JSON(http.StatusOK, movie)
}
