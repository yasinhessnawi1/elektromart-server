package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

// GetOrderItem fetches a single order item by ID provided in the URL.
// It validates the order item data and returns the order item details or an error message if not found or data is invalid.
func GetOrderItem(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var orderItem models.OrderItem

	if err := db.Where("id = ?", id).First(&orderItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}
	c.JSON(http.StatusOK, orderItem)

}

// GetOrderItems retrieves all order items from the database.
// It returns a list of order items in JSON format or an error message if the retrieval fails.
func GetOrderItems(c *gin.Context, db *gorm.DB) {
	orderItems, err := models.GetAllOrderItems(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order items"})
		return
	}
	c.JSON(http.StatusOK, orderItems)
}

func SearchAllOrderItems(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"order_id", "product_id", "quantity", "subtotal"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			switch field {
			case "order_id", "product_id", "quantity":
				if numVal, err := strconv.Atoi(cleanValue); err == nil {
					searchParams[field] = numVal
				}
			case "subtotal":
				if numVal, err := strconv.ParseFloat(cleanValue, 64); err == nil {
					searchParams[field] = numVal
				}
			}
		}
	}

	orderItems, err := models.SearchOrderItem(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order items", "details": err.Error()})
		return
	}

	if len(orderItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No order items found"})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}

// CreateOrderItem handles the creation of a new order item from JSON input.
// It checks the existence of the product, validates input, and persists the new order item in the database.
// Responds with the created order item or an error message.
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

// UpdateOrderItem modifies an existing order item based on the JSON input and the ID provided in the URL.
// It checks product existence, validates the input data, and updates the order item in the database.
// Responds with the updated order item or an error message.
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

// DeleteOrderItem removes an order item from the database based on its ID.
// It handles the deletion process and responds with HTTP 204 No Content on success or an error message if not found or deletion fails.
func DeleteOrderItem(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.OrderItemExists(db, convertedId) {
		fmt.Println("Order item does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order item"})
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
