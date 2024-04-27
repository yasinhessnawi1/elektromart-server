package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetOrders handles the retrieval of all orders from the database.
// It returns a JSON response with a list of orders or an error message if the retrieval fails.
func GetOrders(c *gin.Context, db *gorm.DB) {
	orders, err := models.GetAllOrders(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// CreateOrder handles the creation of a new order from JSON input.
// It validates the input and persists the new order record in the database, returning the created order or an error message.
func CreateOrder(c *gin.Context, db *gorm.DB) {
	var newOrder models.Order
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !tools.CheckDate(newOrder.Order_date) || !tools.CheckStatus(newOrder.Status, 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var order models.Order
	order.Model.ID = uint(tools.GenerateUUID())
	order.User_ID = newOrder.User_ID
	order.Order_date = newOrder.Order_date
	order.Total_amount = newOrder.Total_amount
	order.Status = newOrder.Status
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// GetOrder retrieves a single order by ID from the database.
// It checks the validity of the order data and returns the order details or appropriate error messages.
func GetOrder(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var order models.Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if !tools.CheckDate(order.Order_date) || !tools.CheckStatus(order.Status, 1000) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid order data"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// UpdateOrder handles the updating of an existing order based on the JSON input and the ID provided in the URL.
// It validates the input and updates the order in the database, returning the updated order or an error message.
func UpdateOrder(c *gin.Context, db *gorm.DB) {
	var updatedOrder models.Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !tools.CheckDate(updatedOrder.Order_date) || !tools.CheckStatus(updatedOrder.Status, 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var order models.Order
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	order.User_ID = updatedOrder.User_ID
	order.Order_date = updatedOrder.Order_date
	order.Total_amount = updatedOrder.Total_amount
	order.Status = updatedOrder.Status
	if err := db.Model(&order).Where("id = ?", id).Updates(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedOrder)
}

// DeleteOrder removes an order from the database based on the ID provided in the URL.
// It responds with HTTP 204 No Content on successful deletion or an error message if the order is not found or deletion fails.
func DeleteOrder(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&models.Order{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err := db.Where("id = ?", id).Delete(&models.Order{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// HandleOrderEndpoint routes the HTTP request based on the method to the appropriate handler.
// It supports GET, POST, PUT, and DELETE methods for order management.
func HandleOrderEndpoint(c *gin.Context, db *gorm.DB) {
	switch c.Request.Method {
	case "GET":
		GetOrders(c, db)
	case "POST":
		CreateOrder(c, db)
	case "PUT":
		UpdateOrder(c, db)
	case "DELETE":
		DeleteOrder(c, db)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}
