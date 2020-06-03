package movielight

import (
	"ms/movielight/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
