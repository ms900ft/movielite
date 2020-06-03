package movielight

import (
	"ms/movielight/models"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type FileInput struct {
	FullPath string `json:"fullpath" binding:"required"`
	FileName string `json:"file" `
}

func (s *Service) getFiles(c *gin.Context) {
	db := s.DB

	var files []models.File
	db.Find(&files)
	c.JSON(http.StatusOK, files)

}

func (s *Service) getFile(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var file models.File
	if err := db.Where("id = ?", id).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, file)

}

func (s *Service) addFile(c *gin.Context) {
	db := s.DB
	var input FileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := models.File{FullPath: input.FullPath}

	if err := db.Create(&file).Error; gorm.IsRecordNotFoundError(err) {
		content := gin.H{"error: ": "create" + err.Error()}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
	}
	movie := models.Movie{}
	//movie.FileID = f.ID
	regex := map[string]string{}
	//movie.Title = Translatename(f.FileName)
	movie.Title = Translatename(file.FileName, regex)
	//movie.WatchList = true
	err := movie.GetMeta()
	if err != nil {
		log.Error(err)
	}
	spew.Dump(movie)
	movie.FileID = file.ID
	if err := db.Create(&movie).Error; gorm.IsRecordNotFoundError(err) {
		content := gin.H{"error: ": "create movie" + err.Error()}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
	}
	c.JSON(http.StatusCreated, file)
}

func (s *Service) updateFile(c *gin.Context) {
	db := s.DB
	var input FileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("files binding input: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var file models.File
	if err := db.Where("id = ?", c.Param("id")).First(&file).Error; err != nil {
		log.Errorf("files not found: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input

	if err := db.Model(&file).Updates(input).Error; err != nil {
		log.Errorf("files update error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (s *Service) deleteFile(c *gin.Context) {
	db := s.DB
	var file models.File
	if err := db.Where("id = ?", c.Param("id")).First(&file).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&file)

	c.JSON(http.StatusOK, file)
}
