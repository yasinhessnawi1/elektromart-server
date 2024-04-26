package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetOrderItem(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var orderItem models.OrderItem

	if err := db.Where("id = ?", id).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}
	c.JSON(http.StatusOK, orderItem)

}

func GetOrderItems(c *gin.Context, db *gorm.DB) {
	orderItems, err := models.GetAllOrderItems(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order items"})
		return
	}
	c.JSON(http.StatusOK, orderItems)
}

func CreateOrderItem(c *gin.Context, db *gorm.DB) {
	var newOrderItem models.OrderItem
	if err := c.ShouldBindJSON(&newOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	orderItem := models.OrderItem{
		Order_ID:   newOrderItem.Order_ID,
		Product_ID: newOrderItem.Product_ID,
		Quantity:   newOrderItem.Quantity,
		Subtotal:   newOrderItem.Subtotal,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkOrderItem(orderItem, newOrderItem, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderItem)
}

func UpdateOrderItem(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.OrderItemExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	var updatedOrderItem models.OrderItem
	if err := c.ShouldBindJSON(&updatedOrderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var orderItem models.OrderItem
	if err := db.Where("id = ?", id).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	orderItem.Order_ID = updatedOrderItem.Order_ID
	orderItem.Product_ID = updatedOrderItem.Product_ID
	orderItem.Quantity = updatedOrderItem.Quantity
	orderItem.Subtotal = updatedOrderItem.Subtotal

	if failed, err := checkOrderItem(orderItem, updatedOrderItem, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&orderItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order item", "details": err.Error()})
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

func checkOrderItem(orderItem models.OrderItem, newOrderItem models.OrderItem, db *gorm.DB) (bool, error) {
	switch true {
	case !orderItem.SetOrderID(newOrderItem.Order_ID, db):
		return true, fmt.Errorf("invalid order_id or not existing")
	case !orderItem.SetProductID(newOrderItem.Product_ID, db):
		return true, fmt.Errorf("invalid product_id or not existing")
	case !orderItem.SetQuantity(newOrderItem.Quantity):
		return true, fmt.Errorf("invalid quantity")
	case !orderItem.SetSubtotal(newOrderItem.Subtotal):
		return true, fmt.Errorf("invalid subtotal")
	}
	return false, nil
}
