package movielight

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type genres struct {
	TmdbID int    `json:"tmdb_id"`
	Name   string `json:"name"`
}

func (s *Service) getGenres(c *gin.Context) {
	db := s.DB

	var genres []genres
	if err := db.Select("distinct(name), tmdb_id").Order("name").Find(&genres).Error; err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, genres)
}
