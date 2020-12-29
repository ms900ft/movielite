package models

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	UserID  int64
	Name    string
	IsAdmin bool
	*jwt.StandardClaims
}
