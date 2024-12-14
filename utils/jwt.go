package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID uint, email string) (string, error) {
	secretKey := GetENV().SecretKeyJWT
	claims := jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Mendapatkan SecretKey untuk memvalidasi JWT
	secretKey := GetENV().SecretKeyJWT

	// Mem-parsing token untuk memvalidasi tanda tangan dan mengambil klaim
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Memastikan algoritma tanda tangan adalah HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	// Jika terjadi kesalahan saat parsing
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Memeriksa apakah token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validasi tambahan: memeriksa apakah token sudah kedaluwarsa
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
