package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
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

func CreateOrderItem(c *gin.Context, db *gorm.DB) {
	var newOrderItem models.OrderItemDB
	if err := c.ShouldBindJSON(&newOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var orderItem models.OrderItem
	orderItem.OrderItem_ID = tools.GenerateUUID()
	orderItem.Order_ID = newOrderItem.Order_ID
	orderItem.Product_ID = newOrderItem.Product_ID
	orderItem.Quantity = newOrderItem.Quantity
	orderItem.Subtotal = newOrderItem.Subtotal

	if err := db.Create(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, orderItem)
}

func UpdateOrderItem(c *gin.Context, db *gorm.DB) {
	var updatedOrderItem models.OrderItemDB
	if err := c.ShouldBindJSON(&updatedOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var orderItem models.OrderItem
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}
	if err := db.Where("id = ?", id).Updates(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderItem)
}

func DeleteOrderItem(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&models.OrderItem{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}
	if err := db.Where("id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
