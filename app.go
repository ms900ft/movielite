package movielite

import (
	"fmt"
	"net/http"

	"github.com/gammazero/workerpool"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"

	"github.com/ms900ft/movielite/models"
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
	WorkerPool *workerpool.WorkerPool
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
	if a.Config.SQLDebug {
		a.DB.LogMode(true)
	}
	if err != nil {
		log.Fatal(err)
	}
	tmdbClient := tmdb.Init(tmdb.Config{APIKey: a.Config.TMDBApiKey})
	// pool to slow down tmpdb image requests
	a.WorkerPool = workerpool.New(1)
	models.HttpClient = &http.Client{}
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
	//a.Router.Use(CORSMiddleware())

	a.Router.Use(gin.Recovery())

	a.initializeRoutes()
}

//Run mal sehen
func (a *Service) Run() error {
	p := fmt.Sprintf(":%d", a.Config.Port)
	log.Debug("running on port: " + p)
	err := http.ListenAndServe(p, a.Router)
	return err
}
