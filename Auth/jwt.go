package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("SUPER-SECRET")

type JWTClaim struct {
	User_Name string `json:"username"`
	jwt.StandardClaims
}

func Trim(token string) (string, error) {
	parts := strings.SplitN(token, "Bearer ", 2)
	return Authentication(parts[1])
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		User_Name: username,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func Authentication(SignedStringtoken string) (string, error) {
	token, err := jwt.ParseWithClaims(SignedStringtoken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*JWTClaim)
	username := claims.User_Name
	if !ok {
		err = errors.New("couldn't parse claims")
		return "", err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return "", err
	}
	return username, nil
}
