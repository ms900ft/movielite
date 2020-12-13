package movielight

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
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
