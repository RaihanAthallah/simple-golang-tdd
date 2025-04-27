package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// You can replace these with a config/env loader later
var (
	accessSecret  = []byte(os.Getenv("ACCESS_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
)

func GenerateAccessToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(15 * time.Minute).Unix(), // expires in 15 minutes
		"iat": time.Now().Unix(),                       // issued at
		"iss": "simple-golang-tdd",                     // issuer
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // expires in 7 days
		"iat": time.Now().Unix(),
		"iss": "simple-golang-tdd",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ValidateToken checks the validity of a JWT token and extracts claims
func ValidateToken(tokenString string, tokenType string) (jwt.MapClaims, error) {
	var secret []byte
	if tokenType == "access" {
		secret = accessSecret
	} else if tokenType == "refresh" {
		secret = refreshSecret
	} else {
		return nil, errors.New("invalid token type")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

func NewAccessToken(refreshToken string) (string, error) {
	// Step 1: Validate the refresh token
	claims, err := ValidateToken(refreshToken, "refresh")
	if err != nil {
		return "", err
	}

	// Step 2: Extract user ID from the token claims
	userID, ok := claims["sub"].(string) // No need to dereference
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	// Step 3: Generate a new access token for the user
	return GenerateAccessToken(userID)
}
