package tools

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var mySigningKey = []byte("secret")

// GenerateToken generates a JWT token

// TokenAuthMiddleware is the middleware for JWT authentication
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Strip 'Bearer ' prefix if it exists
		if len(tokenString) > 7 && strings.ToUpper(tokenString[0:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			c.Abort()
			return
		}
	}
}

func GenerateTokenWithClaims(username string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT token using HMAC with SHA-256
	claims := token.Claims.(jwt.MapClaims)   // Use MapClaims for easy map-like syntax with claims

	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expiration set to 72 hours from now

	tokenString, err := token.SignedString(mySigningKey) // Sign the token with our secret key
	return tokenString, err
}
