package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("SUPER-SECRET")

type jwtClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Trim extracts the token from the "Bearer" prefix.
func Trim(token string) (string, error) {
	parts := strings.SplitN(token, "Bearer ", 2)
	return Authentication(parts[1])
}

// GenerateJWT generates a JWT token for the given username.
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwtClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Authentication verifies and extracts the username from a JWT token.
func Authentication(signedStringToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedStringToken, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwtClaim)
	username := claims.Username

	if !ok {
		err = errors.New("couldn't parse claims")
		return "", err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token expired")
		return "", err
	}

	return username, nil
}
