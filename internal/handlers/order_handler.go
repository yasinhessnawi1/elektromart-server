package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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
	var updatedOrder models.Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
