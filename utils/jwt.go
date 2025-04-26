package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// You can replace these with a config/env loader later
var (
	accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
)

func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(), // expires in 15 minutes
		"iat": time.Now().Unix(),                       // issued at
		"iss": "your-app-name",                         // issuer
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // expires in 7 days
		"iat": time.Now().Unix(),
		"iss": "your-app-name",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}
