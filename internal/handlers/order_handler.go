package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetOrders(c *gin.Context, db *gorm.DB) {
	orders, err := models.GetAllOrders(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func CreateOrder(c *gin.Context, db *gorm.DB) {
	var newOrder models.OrderDB
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func GetOrder(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var order models.Order
	if err := db.Where("id = ?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)

}

func UpdateOrder(c *gin.Context, db *gorm.DB) {
	var updatedOrder models.OrderDB
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
