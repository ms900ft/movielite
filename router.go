package movielight

import (
	"github.com/gin-gonic/contrib/expvar"
)

func (a *Service) initializeRoutes() {
	a.Router.GET("/file", a.getFiles)

	a.Router.POST("/file", a.addFile)
	a.Router.GET("/file/:id", a.getFile)
	a.Router.PUT("/file/:id", a.updateFile)
	// a.Router.PUT("/file/:id/move/:dir", moveFile) //update file
	a.Router.DELETE("/file/:id", a.deleteFile)
	a.Router.GET("/file/:id/download", a.downloadFile)
	a.Router.GET("/file/:id/download/:name", a.downloadFile) //name im pfad

	a.Router.GET("/movie", a.getMovies)
	a.Router.GET("/movie/:id", a.getMovie)
	// a.Router.PUT("/movie/:id", updateMovie)
	a.Router.POST("/movie", a.createMovie)
	a.Router.DELETE("/movie/:id", a.deleteMovie)
	// a.Router.GET("/movieMeta/:metaid", getMovieMeta)
	// a.Router.PUT("/movie/:id/addMeta/:metaid", addMeta) //update movie
	// a.Router.PUT("/movie/:id/play", playMovie)
	a.Router.GET("/movie/:id/images", a.getMovieImages)

	a.Router.GET("/genre", a.getGenres)
	a.Router.GET("/country", a.getCountries)
	a.Router.GET("/targets", a.getTargets)

	a.Router.GET("/images/:size/:image", a.getImage)
	//	staticDir := viper.GetString("Frontend.Path")
	//movie2Dir := viper.GetString("Frontend.Path2")
	//	a.Router.Use(favicon.New(staticDir + "/favicon.ico"))
	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// a.Router.Static("/html", staticDir)
	// a.Router.StaticFS("/movie2", statikFS)
	a.Router.GET("/debug/vars", expvar.Handler())
}
