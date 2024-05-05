package tools

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// mySigningKey is a secret key used for signing JWT tokens
var mySigningKey = []byte("secret")

// GenerateToken generates a JWT token

// TokenAuthMiddleware is the middleware for JWT authentication
// It checks the Authorization header for a valid JWT token
// If the token is valid, it sets the username in the request context and calls the next handler
// If the token is invalid, it returns a 401 Unauthorized response
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

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

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

// TokenService this interface make the testing service is more easy
// GenerateTokenWithClaims generate a JWT token with the given username and role
type TokenService interface {
	GenerateTokenWithClaims(username, role string) (string, error)
}

// JWTTokenService is a service for generating JWT tokens
type JWTTokenService struct{}

// GenerateTokenWithClaims generates a JWT token with the given username and role
// It returns the token string and an error if the token generation fails
// The token is signed using HMAC with SHA-256 and has an expiration time of 72 hours
func (service *JWTTokenService) GenerateTokenWithClaims(username string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Create a new JWT token using HMAC with SHA-256
	claims := token.Claims.(jwt.MapClaims)   // Use MapClaims for easy map-like syntax with claims

	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token expiration set to 72 hours from now

	tokenString, err := token.SignedString(mySigningKey) // Sign the token with our secret key
	return tokenString, err
}
