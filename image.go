package movielight

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"

	"ms/movielight/models"
)

func (s *Service) getMovieImages(c *gin.Context) {
	mid, _ := strconv.Atoi(c.Param("id"))
	//mi := MovieInfo{ID: mid}
	images, err := s.getTMDBMovieImages(mid)
	//	var mj JSONB
	//	b, err = json.Marshal(images)
	if err != nil {
		log.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, images)
}

func (s *Service) getImage(c *gin.Context) {
	imageURL := s.Config.TMDBImageURL
	imageDir := s.Config.TMDBImageDir
	size := c.Param("size")
	image := c.Param("image")
	url := fmt.Sprintf("%s/%s/%s", imageURL, size, image)
	imagedir := fmt.Sprintf("%s/%s/%s/%s", imageDir, size, image[0:1], image[1:2])
	imagepath := fmt.Sprintf("%s/%s", imagedir, image)
	//check cache
	if _, err := os.Stat(imagepath); os.IsNotExist(err) {
		client := models.HttpClient
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error(err)
		}

		response, err := client.Do(req)
		log.Debug("Getting image from tmdb")
		if err != nil {
			log.Print(fmt.Sprintf("Can't get %s: %s", url, err))
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Error(err)
			}
			bodyString := string(body)
			log.Errorf("get image status %d: %s", response.StatusCode, bodyString)
			c.String(http.StatusNotFound, bodyString)
			return
		}
		err = os.MkdirAll(imagedir, os.ModePerm)
		if err != nil {
			log.Warnf("creating directory %s", err)
		}

		file, err := os.Create(imagepath)
		if err != nil {
			log.Print(fmt.Sprintf("Can't get %s: %s", url, err))
		}
		defer file.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Error(err)
		}

		log.Debugf("save %s to %s", url, imagedir)
	} else {
		c.Header("X-cache", "HIT")
		log.Debugf("getting image from cache %s", url)
	}

	c.File(imagepath)
}

func (s *Service) getTMDBMovieImages(id int) (*tmdb.MovieImages, error) {
	var options = make(map[string]string)

	// options["append_to_response"] = "credits"
	//options["language"] = DEFLANG
	res, err := s.TMDBClient.GetMovieImages(id, options)
	if err != nil {
		log.Error(err)
		return res, err
	}
	//err = preFetchImages(res)
	return res, nil
}
