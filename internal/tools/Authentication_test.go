package tools

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGenerateToken tests generating a JWT token
// It creates a new instance of the JWTTokenService and calls the GenerateToken function
// It then checks if the token was generated successfully
// Finally, it checks if the token is valid
func TestGenerateTokenWithClaims(t *testing.T) {
	tokenService := JWTTokenService{}

	username := "testuser"
	role := "admin"
	tokenStr, err := tokenService.GenerateTokenWithClaims(username, role)
	assert.Nil(t, err)

	// check the claims
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	assert.Nil(t, err)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		assert.Equal(t, username, claims["username"])
		assert.Equal(t, role, claims["role"])
	} else {
		t.Fail()
	}
}

// TestTokenAuthMiddlewareWithValidToken tests the TokenAuthMiddleware with a valid token
// It creates a new instance of the gin router and sets up a protected route
// It then creates a valid token and sends a request to the protected route with the token
// Finally, it checks if the response status code is 200
func TestTokenAuthMiddlewareWithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TokenAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "access granted"})
	})

	// Valid token
	tokenService := JWTTokenService{}
	token, _ := tokenService.GenerateTokenWithClaims("testuser", "user")

	// Create a test request with the valid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestTokenAuthMiddlewareWithInValidToken tests the TokenAuthMiddleware with an invalid token
// It creates a new instance of the gin router and sets up a protected route
// It then sends a request to the protected route with an invalid token
// Finally, it checks if the response status code is 401
func TestTokenAuthMiddlewareWithInValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(TokenAuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "access granted"})
	})

	// Create a test request with the invalid token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidToken")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
