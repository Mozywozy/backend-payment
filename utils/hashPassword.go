package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte("key")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func ValidateToken(tokenString string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}