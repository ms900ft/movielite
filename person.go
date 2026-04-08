package movielite

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
)

// getPerson godoc
// @Summary Get a person by ID
// @Description get a person by TMDB ID
// @Tags persons
// @Produce  json
// @Param id path int true "Person ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/person/{id} [get]
func (s *Service) getPerson(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid person ID"})
		return
	}

	// First check if we have it in our database
	var crew models.Crew
	if err := s.DB.Where("person_id = ?", id).First(&crew).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"ID": crew.PersonID, "Name": crew.Name})
		return
	}
	var cast models.Cast
	if err := s.DB.Where("person_id = ?", id).First(&cast).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"ID": cast.PersonID, "Name": cast.Name})
		return
	}

	// Not in DB, try to fetch from TMDB
	lang := "en-US"
	tmdbURL := "https://api.themoviedb.org/3/person/" + idStr + "?api_key=" + s.Config.TMDBApiKey + "&language=" + lang

	resp, err := http.Get(tmdbURL)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch from TMDB"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	var tmdbPerson map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tmdbPerson); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse TMDB response"})
		return
	}

	name, _ := tmdbPerson["name"].(string)
	c.JSON(http.StatusOK, gin.H{"ID": id, "Name": name})
}
