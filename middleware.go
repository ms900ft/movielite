package movielite

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("Username %s ", username)

	//db := c.MustGet("DB").(*sql.DB)
	db := s.DB

	var user models.User
	if err := db.Where("user_name  = ?", username).First(&user).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		log.Debugf("no user found %s", err)
	}

	s.User = &user
	//err = user.get(db)
	//c.Set("username", username)
	//c.Set("user", user)
	// if err != nil {
	// 	log.Debugf("no user found %s", err)
	// }

	// Pass on to the next-in-chain
	c.Next()
}

func (s *Service) JwtVerify(c *gin.Context) {

	var header = c.Request.Header.Get("authorization") //Grab the token from the header

	reqToken := strings.TrimSpace(header)
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	if reqToken == "" {
		//Token is missing, returns with error code 403 Unauthorized
		//w.WriteHeader(http.StatusForbidden)
		//json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
		c.JSON(http.StatusForbidden, gin.H{"Message": "Missing auth token"})
		c.Abort()
		return
	}
	tk := &models.Token{}

	_, err := jwt.ParseWithClaims(reqToken, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"Message": err.Error()})
		c.Abort()
		return
	}
	var user models.User
	db := s.DB
	if err := db.Where("user_name  = ?", tk.Name).First(&user).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		log.Debugf("no user found %s", err)
	}
	s.User = &user
	// ctx := context.WithValue(r.Context(), "user", tk)
	// next.ServeHTTP(w, r.WithContext(ctx))
	c.Next()
}
