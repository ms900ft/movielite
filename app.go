package movielight

import (
	"fmt"
	"ms/movielight/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

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
	Router *gin.Engine
	DB     *gorm.DB
	User   *models.User
	Config *Config
	//Config Config
}

type meta struct {
	Total int64 `json:"total"`
}

//Initialize mal sehen
func (a *Service) Initialize(dbname string) {

	var err error
	a.DB = models.ConnectDataBase(dbname) // new
	if err != nil {
		log.Fatal(err)
	}

	if a.Config.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode)
		log.SetLevel(log.InfoLevel)
	} else {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	}

	a.Router = gin.Default()
	a.Router.Use(CORSMiddleware())
	a.Router.Use(a.UserMiddleWare)
	a.Router.Use(gin.Recovery())

	a.initializeRoutes()
}

//Run mal sehen
func (a *Service) Run(addr string) error {
	p := fmt.Sprintf(":%d", a.Config.Port)
	err := http.ListenAndServe(p, a.Router)
	return err
}
