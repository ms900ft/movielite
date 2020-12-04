package movielight

import (
	"net/http"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"ms/movielight/models"
)

func TestMain(m *testing.M) {
	S = Setup()
	S.Initialize()
	S.TMDBClient = &MockTMDBClient{}
	file := models.File{FullPath: "/test/kehraus.mp4"}
	if err := file.Create(S.DB, S.TMDBClient); err != nil {
		log.Fatal("can't create first movie " + err.Error())
	}
	user := models.User{UserName: "test"}
	if err := S.DB.Create(&user).Error; err != nil {
		log.Fatal("can't create first movie " + err.Error())
	}
	code := m.Run()

	//shutdown()
	os.Exit(code)
}

func TestCreateFile_WRONG_DATA(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Post("/api/file").
			JSON(`{"filename":"kehraus.mp4"}`). // request
			Expect(t).
			Status(http.StatusBadRequest).
			End()
}
func TestCreateFile(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Post("/api/file").
			JSON(`{"fullpath":"/test/Paterson.mp4"}`). // request
			Expect(t).
			Assert(jsonpath.Present(`$.ID`)).
			Assert(jsonpath.Contains(`$.FileName`, "Paterson.mp4")).
			Status(http.StatusCreated).
			End()
}

func TestGetFile(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/file/1").
			Expect(t).
			Assert(jsonpath.Present(`$.ID`)).
			Assert(jsonpath.Contains(`$.FileName`, "kehraus.mp4")).
			Status(http.StatusOK).
			End()
}
func TestGetFile_NotFound(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/file/111111").
			Expect(t).
			Status(http.StatusNotFound).
			End()
}

func TestGetFiles(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Get("/api/file").
			Expect(t).
			Assert(jsonpath.GreaterThan(`$`, 1)).
			Status(http.StatusOK).
			End()
}

func TestUpdateFile(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/file/1").
			JSON(`{"fullpath":"/test2/kehraus.mp4"}`). // request
			Expect(t).
			Assert(jsonpath.Present(`$.ID`)).
			Assert(jsonpath.Contains(`$.FileName`, "kehraus.mp4")).
			Assert(jsonpath.Contains(`$.FullPath`, "/test2/kehraus.mp4")).
			Status(http.StatusOK).
			End()
}

func TestUpdateFile_NOT_FOUND(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/file/1111111").
			JSON(`{"fullpath":"/test2/kehraus.mp4"}`). // request
			Expect(t).
			Status(http.StatusNotFound).
			End()
}

func TestUpdateFile_WRONG_DATA(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Put("/api/file/1").
			JSON(`{"filename":"kehraus.mp4"}`). // request
			Expect(t).
			Status(http.StatusBadRequest).
			End()
}

func TestDeleteFile(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Delete("/api/file/1").
			Expect(t).
			Assert(jsonpath.Present(`$.ID`)).
			Assert(jsonpath.Contains(`$.FileName`, "kehraus.mp4")).
			Status(http.StatusOK).
			End()
}

func TestDeleteFile_NOT_FOUND(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Delete("/api/file/111111").
			Expect(t).
			Status(http.StatusNotFound).
			End()
}

// func TestAddFile(t *testing.T) {

// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()
// 	router.POST("/file", S.addFile)
// 	tests := []test{
// 		{
// 			input: []byte(`{"fullpath":"/test/home.mp4"}`),
// 			name:  "add file", want: want{
// 				code: 201, id: 1, filename: "home.mp4",
// 			},
// 		},
// 		{
// 			input: []byte(`{"fullpath":"/test/home.mp4"}`),
// 			name:  "add again file", want: want{
// 				code: 400, id: 0, filename: "",
// 			},
// 		},
// 		{
// 			input: []byte(`{}`),
// 			name:  "empty request", want: want{
// 				code: 400, id: 0, filename: "",
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		r, _ := http.NewRequest("POST", "/files", bytes.NewBuffer(tc.input))
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, r)
// 		resp := w.Result()
// 		var dst struct {
// 			ID       int
// 			Filename string
// 		}
// 		defer resp.Body.Close()
// 		if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
// 			t.Fatal(err)
// 		}

// 		assert.Equal(t, tc.want.code, resp.StatusCode, tc.name)
// 		if dst.ID != tc.want.id {
// 			t.Fatalf("expected 'ID', got '%d'", dst.ID)
// 		}
// 		if dst.Filename != tc.want.filename {
// 			t.Fatalf("expected '%s', got '%s'", tc.want.filename, dst.Filename)
// 		}
// 	}
// }
