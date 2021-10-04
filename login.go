package movielite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ms900ft/movielite/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type token struct {
	UserName string `json:"user_name"`
	Token    string `json:"token"`
	IsAdmin  bool   `json:"is_admin"`
}

type input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	InvalidCredentials = fmt.Errorf("invalid login credentials. Please try again")
	NoSecretFound      = fmt.Errorf("no secret found")
)

func (s *Service) login(c *gin.Context) {
	in := input{}
	err := json.NewDecoder(c.Request.Body).Decode(&in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := s.FindOne(in.Username, in.Password)
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
		log.Errorf("user %s not found", username)
		return resp, InvalidCredentials
	}
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil { //Password does not match!
		log.Debug(err)
		return resp, InvalidCredentials
	}

	tk := &models.Token{
		UserID:  user.ID,
		Name:    user.UserName,
		IsAdmin: user.IsAdmin,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	if s.Config.Secret == "" {
		return resp, NoSecretFound
	}
	tokenString, err := token.SignedString([]byte(s.Config.Secret))
	if err != nil {
		fmt.Println(err)
		return resp, err
	}

	resp.Token = tokenString
	resp.UserName = username
	resp.IsAdmin = user.IsAdmin
	return resp, nil
}
