package movielight

import (
	"fmt"
	"ms/movielight/models"

	"github.com/jinzhu/gorm"
	"github.com/ryanbradynd05/go-tmdb"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	DEFLANG   = "de-DE"
	WATCHLIST = "watchlist"
)

//Service mal sehen
type Service struct {
	Router     *gin.Engine
	DB         *gorm.DB
	User       *models.User
	Config     *Config
	TMDBClient models.TMDBClients
	//Config Config
}

type meta struct {
	Total int64 `json:"total"`
}

//Initialize mal sehen
func (a *Service) Initialize() {
	var err error
	a.DB = models.ConnectDataBase(a.Config.DataBase) // new
	if err != nil {
		log.Fatal(err)
	}
	tmdbClient := tmdb.Init(tmdb.Config{APIKey: a.Config.TMDBApiKey})

	//sa.TMDBClient = new(models.TMDBClient)
	a.TMDBClient = tmdbClient
	if a.Config.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.InfoLevel)
	} else {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	}

	a.Router = gin.Default()
	a.Router.Use(CORSMiddleware())
	//	a.Router.Use(a.UserMiddleWare)
	a.Router.Use(gin.Recovery())

	a.initializeRoutes()
}

//Run mal sehen
func (a *Service) Run() error {
	p := fmt.Sprintf(":%d", a.Config.Port)
	err := http.ListenAndServe(p, a.Router)
	return err
}
