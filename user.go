package movielight

import (
	"ms/movielight/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type UserInput struct {
	UserName string `json:"username" binding:"required"`
}

func (s *Service) getUsers(c *gin.Context) {
	db := s.DB

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, users)

}

func (s *Service) getUser(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, user)

}

func (s *Service) createUser(c *gin.Context) {
	db := s.DB
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{UserName: input.UserName}

	if err := db.Create(&user).Error; gorm.IsRecordNotFoundError(err) {
		content := gin.H{"error: ": "create" + err.Error()}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
	}

	c.JSON(http.StatusCreated, user)
}

func (s *Service) updateUser(c *gin.Context) {
	db := s.DB
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("user binding input: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		log.Errorf("user not found: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input

	if err := db.Model(&user).Updates(input).Error; err != nil {
		log.Errorf("user update error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (s *Service) deleteUser(c *gin.Context) {
	db := s.DB
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&user)

	c.JSON(http.StatusOK, user)
}
