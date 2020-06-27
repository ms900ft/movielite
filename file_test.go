package movielight

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ms/movielight/models"
	"net/http"
	"os"
	"testing"

	"github.com/prometheus/common/log"
	"github.com/ryanbradynd05/go-tmdb"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type test struct {
	input []byte
	name  string
	want  want
}

type want struct {
	code     int
	id       int
	filename string
}
type MockTMDBClient struct{}

var S Service

func (m *MockTMDBClient) SearchMovie(s string, opts map[string]string) (*tmdb.MovieSearchResults, error) {
	log.Error("xxxxxxx search")
	filename := fmt.Sprintf("testdata/search_%s.json", s)
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
	var search tmdb.MovieSearchResults
	json.Unmarshal(byte, &search)
	//spew.Dump(search)
	return &search, nil
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

func TestMain(m *testing.M) {
	S = setup()
	S.Initialize()
	S.TMDBClient = &MockTMDBClient{}
	code := m.Run()
	//shutdown()
	os.Exit(code)
}

func setup() Service {

	c := Config{}
	c.Mode = "testing"
	s := Service{Config: &c}
	db := models.ConnectDataBase(":memory:")
	s.DB = db
	user := models.User{UserName: "marc"}
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	return s
}

func TestCreateGetFile(t *testing.T) {
	apitest.New(). // configuration
			Handler(S.Router).
			Post("/file").
			JSON(`{"fullpath":"/test/Paterson.mp4"}`). // request
			Expect(t).
			Assert(jsonpath.Present(`$.ID`)).
			Assert(jsonpath.Contains(`$.FileName`, "Paterson.mp4")).
			Status(http.StatusCreated).
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
