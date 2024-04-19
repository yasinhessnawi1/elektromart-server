package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetCategories(c *gin.Context, db *gorm.DB) {
	categories, err := models.GetAllCategories(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}
