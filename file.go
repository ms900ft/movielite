package movielight

import (
	"ms/movielight/models"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type FileInput struct {
	FullPath string `json:"fullpath" binding:"required"`
	FileName string `json:"file" `
}

func (s *Service) getFiles(c *gin.Context) {
	db := s.DB

	var files []models.File
	if err := db.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
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
	if err := file.Create(db, s.TMDBClient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//spew.Dump(movie)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&file)

	c.JSON(http.StatusOK, file)
}

func (s *Service) downloadFile(c *gin.Context) {
	db := s.DB
	id := c.Param("id")

	var file models.File
	if err := db.Where("id = ?", id).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+file.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.FullPath)
}
func (s *Service) moveFile(c *gin.Context) {
	db := s.DB
	id := c.Param("id")
	dir := c.Param("dir")
	var file models.File
	if err := db.Where("id = ?", id).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	newpath, err := file.Move(dir)
	if err != nil {
		log.Error("move file:  " + err.Error())
		content := gin.H{"error": "File not found"}
		log.Error(content)
		c.JSON(http.StatusBadRequest, content)
		return
	}
	file.FullPath = newpath
	if err := db.Model(&file).Updates(file).Error; err != nil {
		log.Errorf("files update error: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}
	c.JSON(http.StatusOK, file)
}
