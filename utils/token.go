package utils

import (
	"awesomeProject/model/common"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	InValidToken     = errors.New("invalid token")
)

type JWT struct {
	SignedKey []byte
}

func NewJwt() *JWT {
	return &JWT{
		SignedKey: []byte("斗龙战士"),
	}

}

func (j *JWT) CreateToken(UserName string, UserID int64) (string, error) {

	claims := common.Claims{
		UserId:   UserID,
		UserName: UserName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignedKey)

}

func (j *JWT) ParseToken(tokenString string) (*common.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &common.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SignedKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, InValidToken
			}

		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*common.Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, InValidToken
	} else {
		return nil, InValidToken
	}

}
