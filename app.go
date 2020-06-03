package movielight

import (
	"ms/movielight/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//Service mal sehen
type Service struct {
	Router *gin.Engine
	DB     *gorm.DB
	//Config Config
}

type meta struct {
	Total int64 `json:"total"`
}

//
//Query movie query parameter
type Query struct {
	Orderby     string `form:"orderby"`
	Qtitel      string `form:"title"`
	Alpha       string `form:"alpha"`
	Genre       int64  `form:"genre"`
	Crew        int64  `form:"crew"`
	Person      int64  `form:"person"`
	Cast        int64  `form:"cast"`
	Country     string `form:"country"`
	LastScanned string `form:"last_scanned"`
}

//Initialize mal sehen
func (a *Service) Initialize(user, host, password, port, dbname string) {

	var err error
	a.DB = models.ConnectDataBase() // new
	if err != nil {
		log.Fatal(err)
	}

	mode := gin.Mode()
	if mode == "release" {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	a.Router = gin.Default()
	a.Router.Use(CORSMiddleware())
	//a.Router.Use(Database())
	//a.Router.Use(UserMiddleWare)

	//	a.DB = c.MustGet("DB").(*sql.DB)
	a.initializeRoutes()
}

//Run mal sehen
func (a *Service) Run(addr string) error {
	err := http.ListenAndServe(":8000", a.Router)
	return err
}
