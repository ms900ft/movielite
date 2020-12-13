package movielight

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"

	"ms/movielight/models"
)

type MockTMDBClient struct{}

var S Service

func (m *MockTMDBClient) SearchMovie(s string, opts map[string]string) (*tmdb.MovieSearchResults, error) {
	log.Error("xxxxxxx search")
	filename := fmt.Sprintf("testdata/search_%s.json", s)
	jsonFile, err := os.Open(filename)
	var search tmdb.MovieSearchResults
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	fmt.Println("Successfully Opened " + filename)
	defer jsonFile.Close()
	byte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	json.Unmarshal(byte, &search)
	//spew.Dump(search)
	return &search, nil
}
func (m *MockTMDBClient) GetMovieImages(id int, opts map[string]string) (*tmdb.MovieImages, error) {
	log.Debug("mock images")
	filename := fmt.Sprintf("testdata/images_%d.json", id)
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Opened " + filename)
	defer jsonFile.Close()
	byte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var images tmdb.MovieImages
	json.Unmarshal(byte, &images)
	return &images, nil
}

func (m *MockTMDBClient) GetMovieInfo(id int, opts map[string]string) (*tmdb.Movie, error) {
	log.Error("xxxxxxx search")
	filename := fmt.Sprintf("testdata/movie_%d.json", id)
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Opened " + filename)
	defer jsonFile.Close()
	byte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var movie tmdb.Movie
	json.Unmarshal(byte, &movie)
	return &movie, nil
}

func Setup() Service {

	c := Config{}
	c.Mode = "testing"
	c.TargetDir = "./testdata"
	s := Service{Config: &c}
	s.TMDBClient = &MockTMDBClient{}
	db := models.ConnectDataBase(":memory:")
	s.DB = db
	user := models.User{UserName: "marc"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	s.User = &user
	return s
}
