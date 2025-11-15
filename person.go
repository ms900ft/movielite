package movielite

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
)

// getPerson godoc
// @Summary Get a person by ID
// @Description get a person by TMDB ID
// @Tags persons
// @Produce  json
// @Param id path int true "Person ID"
// @Success 200 {object} models.Person
// @Router /api/person/{id} [get]
func (s *Service) getPerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	var crew models.Crew
	if err := s.DB.Where("person_id = ?", id).First(&crew).Error; err == nil {
		c.JSON(http.StatusOK, crew)
		return
	}
	var cast models.Cast
	if err := s.DB.Where("person_id = ?", id).First(&cast).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// Fetch from TMDB and store

		} else {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	c.JSON(http.StatusOK, cast)
}
