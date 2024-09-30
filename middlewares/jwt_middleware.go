package middleware

import (
	"net/http"
	"strings"
	"time"

	"payment-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		autHeader := c.GetHeader("Authorization")
		if autHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Authorization header is missing"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(autHeader, " ")

		useID, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("userId", useID)

		c.Next()
	}
}

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(utils.JwtKey))
}