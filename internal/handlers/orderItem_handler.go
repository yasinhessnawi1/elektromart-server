package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetOrderItems(c *gin.Context, db *gorm.DB) {
	orderItems, err := models.GetAllOrderItems(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order items"})
		return
	}
	c.JSON(http.StatusOK, orderItems)
}
