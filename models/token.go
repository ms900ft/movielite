package models

import jwt "github.com/golang-jwt/jwt/v4"

//Token struct declaration
type Token struct {
	UserID  int64
	Name    string
	IsAdmin bool
	*jwt.StandardClaims
}
