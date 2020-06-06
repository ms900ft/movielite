package movielight

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type countries struct {
	Iso3166_1 string `json:"iso_id"`
	Name      string `json:"name"`
}

func (s *Service) getCountries(c *gin.Context) {
	db := s.DB

	var countries []countries
	if err := db.Table("production_countries").Select("distinct(name), iso3166_1").Order("name").Find(&countries).Error; err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, countries)
}
