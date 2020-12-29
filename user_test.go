package movielite

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestUserGet(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/user").
			Expect(t).
		//Body(`[{"UserName":"admin", "id":1}]`).
		Assert(jsonpath.GreaterThan(`$`, 2)).
		Status(http.StatusOK).
		End()
}

func TestUserGetOne(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/user/1").
			Expect(t).
			Body(`{"UserName":"admin", "id":1, "is_admin": true}`).
			Status(http.StatusOK).
			End()
}
