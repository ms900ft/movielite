package movielight

import (
	"ms/movielight/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

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
	err := c.Bind(&q)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	db := s.DB
	var movies []models.Movie
	var count int64
	count = 2300
	if err := db.Set("gorm:auto_preload", true).Model(&models.Movie{}).
		Order("created_at DESC").Count(&count).Offset(q.Offset).Limit(q.Limit).Find(&movies).
		Error; err != nil {
		log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
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

func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode_sign = ?", "C")
}
