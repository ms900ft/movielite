package movielight

import (
	"net/http"
	"os"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestTargetsGet(t *testing.T) {
	S.Config.TargetDir = t.TempDir()
	err := os.Mkdir(S.Config.TargetDir+"/testdir", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/targets").
			Expect(t).
			Assert(jsonpath.GreaterThan(`$`, 1)).
			Body(`[{"name": "testdir"}]`).
			Status(http.StatusOK).
			End()
}

func TestTargetsGetWrongDir(t *testing.T) {
	S.Config.TargetDir = "/tmp/hsashsahjsajsjsa"
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/targets").
			Expect(t).
			Body(`{"error":"open /tmp/hsashsahjsajsjsa: no such file or directory"}`).
			Status(http.StatusInternalServerError).
			End()
}
