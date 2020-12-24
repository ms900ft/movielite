package movielite

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	pass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
	}

	user := models.User{UserName: input.UserName, Password: string(pass)}

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
	pass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
	}
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		log.Errorf("user not found: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	input.Password = string(pass)
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
