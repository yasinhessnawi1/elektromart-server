package handlers

import (
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostLogin(c *gin.Context) {
	var loginCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect parameters"})
		return
	}

	// Dummy check for username and password (replace with actual DB check)
	if loginCredentials.Username == "admin" && loginCredentials.Password == "admin" {
		tokenString, err := tools.GenerateToken(loginCredentials.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
	}
}
