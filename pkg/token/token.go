package token

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ms900ft/movielite/models"
)

func AdminToken(secret string) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 1000000).Unix()
	tk := &models.Token{
		Name: "admin",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	if secret == "" {
		return "", errors.New("no secret found")
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenString, err
}
