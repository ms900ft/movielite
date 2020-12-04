package movielight

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestGenresGet(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/genre").
			Expect(t).
			Assert(jsonpath.GreaterThan(`$`, 1)).
			Status(http.StatusOK).
			End()
}
