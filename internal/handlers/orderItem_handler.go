package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetOrderItems retrieves all order items from the database.
// It returns a list of order items in JSON format or an error message if the retrieval fails.
func GetOrderItems(c *gin.Context, db *gorm.DB) {
	orderItems, err := models.GetAllOrderItems(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order items", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderItems)
}

// GetOrderItem fetches a single order item by ID provided in the URL.
// It validates the order item data and returns the order item details or an error message if not found or data is invalid.
func GetOrderItem(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var orderItem models.OrderItem
	if err := db.Where("id = ?", id).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}
	if !tools.CheckInt(orderItem.Quantity) || !tools.CheckFloat(orderItem.Subtotal) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid order item data"})
		return
	}
	c.JSON(http.StatusOK, orderItem)
}

// CreateOrderItem handles the creation of a new order item from JSON input.
// It checks the existence of the product, validates input, and persists the new order item in the database.
// Responds with the created order item or an error message.
func CreateOrderItem(c *gin.Context, db *gorm.DB) {
	var newOrderItem models.OrderItem
	if err := c.ShouldBindJSON(&newOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !models.ProductExists(db, newOrderItem.Product_ID) || !tools.CheckInt(newOrderItem.Quantity) || !tools.CheckFloat(newOrderItem.Subtotal) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var orderItem models.OrderItem
	orderItem.Model.ID = uint(tools.GenerateUUID())
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

// UpdateOrderItem modifies an existing order item based on the JSON input and the ID provided in the URL.
// It checks product existence, validates the input data, and updates the order item in the database.
// Responds with the updated order item or an error message.
func UpdateOrderItem(c *gin.Context, db *gorm.DB) {
	var updatedOrderItem models.OrderItem
	if err := c.ShouldBindJSON(&updatedOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !models.ProductExists(db, updatedOrderItem.Product_ID) || !tools.CheckInt(updatedOrderItem.Quantity) || !tools.CheckFloat(updatedOrderItem.Subtotal) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var orderItem models.OrderItem
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}
	orderItem.Order_ID = updatedOrderItem.Order_ID
	orderItem.Product_ID = updatedOrderItem.Product_ID
	orderItem.Quantity = updatedOrderItem.Quantity
	orderItem.Subtotal = updatedOrderItem.Subtotal
	if err := db.Where("id = ?", id).Updates(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orderItem)
}

// DeleteOrderItem removes an order item from the database based on its ID.
// It handles the deletion process and responds with HTTP 204 No Content on success or an error message if not found or deletion fails.
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
