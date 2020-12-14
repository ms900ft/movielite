package models

import (
	"net/http"
)

var (
	HttpClient HTTPClients
)

type HTTPClients interface {
	Do(req *http.Request) (*http.Response, error)
}
