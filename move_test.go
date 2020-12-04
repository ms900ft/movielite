package movielight

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"ms/movielight/models"
)

func TestGetMovie(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/movie/1").
			Expect(t).
			Assert(jsonpath.Present(`$.meta`)).
			Assert(jsonpath.Contains(`$.title`, "kehraus")).
			Assert(jsonpath.Contains(`$.meta.Title`, "Kehraus")).
			Status(http.StatusOK).
			End()
}

func TestGetMovies(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/movie").
			Expect(t).
			Assert(jsonpath.GreaterThan(`$`, 1)).
			Status(http.StatusOK).
			End()
}

func TestUpdateMovie(t *testing.T) {
	mo := models.Movie{}

	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/movie/1").
			JSON(mo).
			Expect(t).
			Assert(jsonpath.Contains(`$.error`, "Movie.Title")).
			Status(http.StatusBadRequest).
			End()
	mo.Title = "Patersonx"
	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/movie/1").
			JSON(mo).
			Expect(t).
			Assert(jsonpath.Contains(`$.title`, "Patersonx")).
			Assert(jsonpath.NotPresent(`$.meta`)).
			Status(http.StatusOK).
			End()

	mo.Title = "Paterson"
	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/movie/1").
			JSON(mo).
			Expect(t).
			Assert(jsonpath.Contains(`$.title`, "Paterson")).
			Assert(jsonpath.Present(`$.meta`)).
			Status(http.StatusOK).
			End()

	mo.Title = "Nachtblende"
	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/movie/1").
			JSON(mo).
			Expect(t).
			Assert(jsonpath.Contains(`$.title`, "Nachtblende")).
			Assert(jsonpath.NotPresent(`$.meta`)).
			Assert(jsonpath.Present(`$.multiplechoice`)).
			Status(http.StatusOK).
			End()
}
