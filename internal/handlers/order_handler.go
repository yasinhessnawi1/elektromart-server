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

func GetOrder(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var order models.Order

	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)

}

func GetOrders(c *gin.Context, db *gorm.DB) {
	orders, err := models.GetAllOrders(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

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
