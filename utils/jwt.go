package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID uint, email string) (string, error) {
	secretKey := "your_secret_key" // Replace with a secure key from environment variables
	claims := jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
