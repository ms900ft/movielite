package movielight

import (
	"fmt"
	"ms/movielight/models"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

//CORSMiddleware for gin
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cookie")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

//UserMiddleWare Get username from movieuser cookie and set username in context
func (s *Service) UserMiddleWare(c *gin.Context) {
	username, err := c.Cookie("movieuser")
	if err != nil {
		log.Debug(err)
		username = "marc"
	}

	username = strings.ToLower(username)
	log.Debugf("Usernamex %s ", username)

	//db := c.MustGet("DB").(*sql.DB)
	db := s.DB

	var user models.User
	if err := db.Where("user_name  = ?", username).First(&user).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	s.User = &user
	//err = user.get(db)
	//c.Set("username", username)
	//c.Set("user", user)
	if err != nil {
		log.Error(err)
	}

	// Pass on to the next-in-chain
	c.Next()
}
