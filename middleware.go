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

func (s *Service) IsAdmin(c *gin.Context) {
	if !s.Token.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"Message": "forbidden"})
		c.Abort()
		return
	}
	c.Next()
}

func (s *Service) JwtVerify(c *gin.Context) {

	var header = c.Request.Header.Get("authorization") //Grab the token from the header
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Missing auth header"})
		c.Abort()
		return
	}
	reqToken := strings.TrimSpace(header)
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	if reqToken == "" {
		//Token is missing, returns with error code 403 Unauthorized
		//w.WriteHeader(http.StatusForbidden)
		//json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Missing auth token"})
		c.Abort()
		return
	}
	tk := &models.Token{}

	_, err := jwt.ParseWithClaims(reqToken, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config.Secret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
		c.Abort()
		return
	}
	var user models.User
	db := s.DB
	if err := db.Where("user_name  = ?", tk.Name).First(&user).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		log.Debugf("no user found %s", err)
	}
	s.Token = tk
	//spew.Dump(tk)
	// ctx := context.WithValue(r.Context(), "user", tk)
	// next.ServeHTTP(w, r.WithContext(ctx))
	c.Next()
}
