package jwtutils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID int, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24*time.Hour).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func JWTMidlleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error": "Authorization header is required"})
			return
		}
		tokenParts := strings.Split(authHeader," ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error": "Invalid token format"})
			return
		}

		token, err := jwt.Parse(tokenParts[1], func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil	
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error": "Invalid token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["sub"]) 
		}

		c.Next()
	}
} 