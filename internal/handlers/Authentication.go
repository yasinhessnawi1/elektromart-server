package handlers

import (
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func PostLogin(c *gin.Context, db *gorm.DB) {
	var loginCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Attempt to bind JSON payload to struct
	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect parameters"})
		return
	}

	// Fetch user from database based on username
	user, err := GetUserByUN(loginCredentials.Username, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed", "details": "user not found"})
		return
	}

	// Compare provided password with the hashed password from database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed", "details": "incorrect password"})
		return
	}

	// Generate token with claims
	tokenString, err := tools.GenerateTokenWithClaims(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	// Return the token string in response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
