package models

import "github.com/ryanbradynd05/go-tmdb"

type TMDBClients interface {
	SearchMovie(string, map[string]string) (*tmdb.MovieSearchResults, error)
	GetMovieInfo(int, map[string]string) (*tmdb.Movie, error)
}

//type TMDBClient struct{}
