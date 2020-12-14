package movielight

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"ms/movielight/models"
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
	json := `{"name":"Test Name","full_name":"test full name","owner":{"login": "octocat"}}`
	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/images/160/xxxxxxxx").
			Expect(t).
		// Assert(jsonpath.GreaterThan(`$.Backdrops`, 1)).
		// Assert(jsonpath.GreaterThan(`$.Posters`, 1)).
		// Assert(jsonpath.Equal(`$.ID`, float64(3476))).
		Status(http.StatusOK).
		End()
}
