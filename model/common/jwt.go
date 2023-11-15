package common

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.StandardClaims
	UserId   int64
	UserName string
}
