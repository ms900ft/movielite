package movielite

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
)

type FileInput struct {
	FullPath string `json:"fullpath" binding:"required"`
	FileName string `json:"file" `
}

// getFiles godoc
// @Summary Get a list of files
// @Description Get a list of all files, optionally filtered by filename.
// @Tags files
// @Produce  json
// @Param   f    query     string  false  "filename to filter by"
// @Success 200 {array} models.File
// @Router /api/file [get]
func (s *Service) getFiles(c *gin.Context) {
	db := s.DB

	var files []models.File
	query := c.DefaultQuery("f", "")
	if query != "" {
		if err := db.Where("file_name = ?", query).Find(&files).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	} else {
		if err := db.Find(&files).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
	if len(files) == 0 {
		c.JSON(http.StatusNotFound, files)
		return
	}
	c.JSON(http.StatusOK, files)
}

// getFile godoc
// @Summary Get a single file
// @Description Get a single file by its ID.
// @Tags files
// @Produce  json
// @Param   id    path      int  true  "File ID"
// @Success 200 {object} models.File
// @Router /api/file/{id} [get]
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

// addFile godoc
// @Summary Add a new file
// @Description Adds a new file to the database.
// @Tags files
// @Accept  json
// @Produce  json
// @Param   file    body      FileInput  true  "File to add"
// @Success 201 {object} models.File
// @Router /api/file [post]
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

// updateFile godoc
// @Summary Update a file
// @Description Updates a file's information.
// @Tags files
// @Accept  json
// @Produce  json
// @Param   id      path      int        true  "File ID"
// @Param   file    body      FileInput  true  "File information to update"
// @Success 200 {object} models.File
// @Router /api/file/{id} [put]
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

// deleteFile godoc
// @Summary Delete a file
// @Description Deletes a file by its ID.
// @Tags files
// @Produce  json
// @Param   id    path      int  true  "File ID"
// @Success 200 {object} models.File
// @Router /api/file/{id} [delete]
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

// downloadFile godoc
// @Summary Download a file
// @Description Downloads a file by its ID.
// @Tags files
// @Produce  octet-stream
// @Param   id    path      int  true  "File ID"
// @Param   name  path      string  false  "Optional download filename"
// @Success 200 {file} file
// @Router /api/file/{id}/download [get]
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

// moveFile godoc
// @Summary Move a file
// @Description Moves a file to a new directory.
// @Tags files
// @Produce  json
// @Param   id    path      int     true  "File ID"
// @Param   dir   path      string  true  "Target directory name"
// @Success 200 {object} models.File
// @Router /api/file/{id}/move/{dir} [put]
func (s *Service) moveFile(c *gin.Context) {
	db := s.DB
	id := c.Param("id")
	dir := c.Param("dir")
	var file models.File
	if err := db.Where("id = ?", id).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}
	newDir := filepath.Join(s.Config.TargetDir, dir)
	newpath, err := file.Move(newDir)
	if err != nil {
		log.Error("move file:  " + err.Error())
		content := gin.H{"error": "file unable to move"}
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
