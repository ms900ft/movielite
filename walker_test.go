package movielite

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ms900ft/movielite/models"
)

// type MockConfig Config

func TestWalker_send(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/file", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var pl payload
		err = json.Unmarshal(body, &pl)
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/file", pl.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	walker := &Walker{
		Config: &Config{
			ServerURL: mockServer.URL,
		},
		Token: "test-token",
	}

	err := walker.send("/path/to/file")
	assert.NoError(t, err)
}
func TestWalker_update(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/file/123", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		var file models.File
		err = json.Unmarshal(body, &file)
		assert.NoError(t, err)
		assert.Equal(t, 123, file.ID)
		assert.Equal(t, "/path/to/file", file.FullPath)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	walker := &Walker{
		Config: &Config{
			ServerURL: mockServer.URL,
		},
		Token: "test-token",
	}
	file := models.File{
		// ID:       123,
		FullPath: "/path/to/file",
	}

	err := walker.update(file)
	assert.NoError(t, err)
}
