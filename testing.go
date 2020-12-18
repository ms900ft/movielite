package movielite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"

	"ms/movielite/models"
)

type MockTMDBClient struct{}

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var S Service

var (
	// GetDoFunc fetches the mock client's `Do` func
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockTMDBClient) SearchMovie(s string, opts map[string]string) (*tmdb.MovieSearchResults, error) {
	log.Debug("search movie")
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
	log.Debug("movie info")
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

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func Setup() Service {

	c := Config{}
	c.Mode = "testing"
	c.TargetDir = "./testdata"
	c.TMDBImageDir = "/tmp/test"
	s := Service{Config: &c}
	s.TMDBClient = &MockTMDBClient{}
	models.HttpClient = &http.Client{}
	db := models.ConnectDataBase(":memory:")
	s.DB = db
	user := models.User{UserName: "marc"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	s.User = &user
	return s
}
