package movielite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type token struct {
	UserName string `json:"user_name"`
	Token    string `json:"token"`
}

func (s *Service) login(c *gin.Context) {

	user := &models.User{}
	err := json.NewDecoder(c.Request.Body).Decode(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := s.FindOne(user.UserName, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (s *Service) FindOne(username, password string) (token, error) {
	db := s.DB
	user := &models.User{}
	resp := token{}

	if err := db.Where("user_name = ?", username).First(user).Error; err != nil {
		log.Debug(err)
		return resp, fmt.Errorf("user %s not found", username)
	}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil { //Password does not match!
		log.Debug(err)
		var err = fmt.Errorf("Invalid login credentials. Please try again %s", err)
		return resp, err
	}

	tk := &models.Token{
		UserID: user.ID,
		Name:   user.UserName,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	resp.Token = tokenString
	resp.UserName = username
	return resp, nil
}
