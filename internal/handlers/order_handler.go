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

// GetOrder retrieves a single order by ID from the database.
// It checks the validity of the order data and returns the order details or appropriate error messages.
func GetOrder(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)

}

// GetOrders handles the retrieval of all orders from the database.
// It returns a JSON response with a list of orders or an error message if the retrieval fails.
func GetOrders(c *gin.Context, db *gorm.DB) {
	orders, err := models.GetAllOrders(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// SearchAllOrders retrieves all orders from the database based on the search parameters provided in the query string.
// It responds with a list of orders if successful or an informational message if no orders exist.
// On failure, it returns an HTTP 500 Internal Server Error.
// The search parameters include user_id, order_date, total_amount, and status.
func SearchAllOrders(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"user_id", "order_date", "total_amount", "status"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			switch field {
			case "user_id":
				if numVal, err := strconv.Atoi(cleanValue); err == nil {
					searchParams[field] = numVal
				}
			case "order_date":
				searchParams[field] = cleanValue
			case "total_amount":
				if numVal, err := strconv.ParseFloat(cleanValue, 64); err == nil {
					searchParams[field] = numVal
				}
			case "status":
				searchParams[field] = strings.ToLower(cleanValue)
			default:
				searchParams[field] = cleanValue
			}
		}
	}

	orders, err := models.SearchOrder(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No orders found"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// CreateOrder handles the creation of a new order based on the JSON input.
// It validates the input and creates the order in the database, returning the created order or an error message.
// The user_id, order_date, total_amount, and status fields are validated for correct formatting.
func CreateOrder(c *gin.Context, db *gorm.DB) {
	var newOrder models.Order
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	order := models.Order{
		User_ID:      newOrder.User_ID,
		Order_date:   newOrder.Order_date,
		Total_amount: newOrder.Total_amount,
		Status:       newOrder.Status,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkOrder(order, newOrder, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// UpdateOrder handles the updating of an existing order based on the JSON input and the ID provided in the URL.
// It validates the input and updates the order in the database, returning the updated order or an error message.
// The user_id, order_date, total_amount, and status fields are validated for correct formatting.
func UpdateOrder(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.OrderExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var updatedOrder models.Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var order models.Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.User_ID = updatedOrder.User_ID
	order.Order_date = updatedOrder.Order_date
	order.Total_amount = updatedOrder.Total_amount
	order.Status = updatedOrder.Status

	if failed, err := checkOrder(order, updatedOrder, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// DeleteOrder removes an order from the database based on the ID provided in the URL.
// It responds with HTTP 204 No Content on successful deletion or an error message if the order is not found or deletion fails.
// The order is soft-deleted using the Unscoped method to preserve data integrity.
func DeleteOrder(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.OrderExists(db, convertedId) {
		fmt.Println("Order does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.Order{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting order"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// checkOrder validates the input data for an order and returns an error if the data is invalid.
// It checks the order's user_id, order_date, total_amount, and status fields for correct formatting.
func checkOrder(order models.Order, newOrder models.Order, db *gorm.DB) (bool, error) {
	switch true {
	case !order.SetUserID(newOrder.User_ID, db):
		return true, fmt.Errorf("invalid user_id or not existing")
	case !order.SetOrderDate(newOrder.Order_date):
		return true, fmt.Errorf("order date is not expected")
	case !order.SetTotalAmount(newOrder.Total_amount):
		return true, fmt.Errorf("invalid amount")
	case !order.SetStatus(newOrder.Status):
		return true, fmt.Errorf("payment status is not expected")
	}
	return false, nil
}
