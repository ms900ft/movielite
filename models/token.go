package models

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	UserID int64
	Name   string
	//	Email  string
	*jwt.StandardClaims
}
