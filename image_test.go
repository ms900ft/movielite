package movielite

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"ms/movielite/models"
)

func TestMovieImageGet(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/movie/3476/images").
			Expect(t).
			Assert(jsonpath.GreaterThan(`$.Backdrops`, 1)).
			Assert(jsonpath.GreaterThan(`$.Posters`, 1)).
			Assert(jsonpath.Equal(`$.ID`, float64(3476))).
			Status(http.StatusOK).
			End()
}

func TestImageGet(t *testing.T) {
	models.HttpClient = &MockHttpClient{}
	json := `{"some body"}}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		header := http.Header{}
		header.Add("Content-Type", "image/jpg")
		return &http.Response{
			StatusCode: 200,
			Body:       r,
			Header:     header,
		}, nil
	}
	S.Config.TMDBImageDir = t.TempDir()
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/images/160/xxxxxxxx").
			Expect(t).
			HeaderNotPresent("X-cache").
			Status(http.StatusOK).
			End()
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/images/160/xxxxxxxx").
			Expect(t).
			HeaderPresent("X-cache").
			Status(http.StatusOK).
			End()
}

func TestImageGetNotFound(t *testing.T) {
	models.HttpClient = &MockHttpClient{}
	json := `{"some body"}}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		header := http.Header{}
		return &http.Response{
			StatusCode: 404,
			Body:       r,
			Header:     header,
		}, nil
	}
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/images/160/xxxxxxxx").
			Expect(t).
			Status(http.StatusNotFound).
			End()

}
