package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetBrands(c *gin.Context, db *gorm.DB) {
	brands, err := models.GetAllBrands(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}
