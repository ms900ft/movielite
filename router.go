package movielite

import (
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/contrib/expvar"
	_ "github.com/ms900ft/movielite/statik"
	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
)

func (a *Service) initializeRoutes() {
	a.Router.Use(CORSMiddleware())
	a.Router.POST("/login", a.login)
	api := a.Router.Group("/api")
	//	api.Use(CORSMiddleware())
	if a.Config.UseAuthentication {
		api.Use(a.JwtVerify)
	}
	api.GET("/file", a.getFiles)

	api.POST("/file", a.addFile)
	api.GET("/file/:id", a.getFile)
	api.PUT("/file/:id", a.updateFile)
	api.PUT("/file/:id/move/:dir", a.moveFile) //update file
	api.DELETE("/file/:id", a.deleteFile)
	api.GET("/file/:id/download", a.downloadFile)
	api.GET("/file/:id/download/:name", a.downloadFile) //name im pfad

	api.GET("/movie", a.getMovies)
	api.GET("/movie/:id", a.getMovie)
	api.PUT("/movie/:id", a.updateMovie)
	api.POST("/movie", a.createMovie)
	api.DELETE("/movie/:id", a.deleteMovie)
	api.PUT("/movie/:id/play", a.playMovie)

	api.GET("/user", a.getUsers)
	api.GET("/user/:id", a.getUser)
	api.PUT("/user/:id", a.IsAdmin, a.updateUser)
	api.POST("/user", a.IsAdmin, a.createUser)
	api.DELETE("/user/:id", a.IsAdmin, a.deleteUser)

	// api.GET("/movieMeta/:metaid", getMovieMeta)
	api.PUT("/movie/:id/addMeta/:metaid", a.addMeta) //update movie

	api.GET("/movie/:id/images", a.getMovieImages)

	api.GET("/genre", a.getGenres)
	api.GET("/country", a.getCountries)
	api.GET("/targets", a.getTargets)

	a.Router.GET("/images/:size/:image", a.getImage)
	//staticDir := viper.GetString("Frontend.Path")
	//movie2Dir := viper.GetString("Frontend.Path2")
	//	a.Router.Use(favicon.New(staticDir + "/favicon.ico"))
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	if a.Config.WebDav {
		WebDAV(a.Config.TargetDir, a.Router.Group("/webdav"))
		log.Info("webdav: /webdav/ waiting for connection")
	}
	// a.Router.Static("/html", staticDir)
	a.Router.StaticFS("/movie2", &indexWrapper{statikFS})
	a.Router.GET("/debug/vars", expvar.Handler())
}

type indexWrapper struct {
	assets http.FileSystem
}

func (i *indexWrapper) Open(name string) (http.File, error) {
	ret, err := i.assets.Open(name)
	if !os.IsNotExist(err) || path.Ext(name) != "" {
		return ret, err
	}

	return i.assets.Open("/index.html")
}
