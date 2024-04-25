package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetPayments(c *gin.Context, db *gorm.DB) {
	payments, err := models.GetAllPayments(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving payments"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func GetPayment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var payment models.Payment
	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func CreatePayment(c *gin.Context, db *gorm.DB) {
	var newPayment models.PaymentDB
	if err := c.ShouldBindJSON(&newPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var payment models.Payment
	payment.Model.ID = uint(tools.GenerateUUID())
	payment.Order_ID = newPayment.ID
	payment.Payment_method = newPayment.Payment_method
	payment.Amount = newPayment.Amount
	payment.Payment_date = newPayment.Payment_date
	payment.Status = newPayment.Status
	if err := db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, payment)
}

func UpdatePayment(c *gin.Context, db *gorm.DB) {
	var updatedPayment models.PaymentDB
	if err := c.ShouldBindJSON(&updatedPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var payment models.Payment
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	payment.Order_ID = updatedPayment.Order_ID
	payment.Payment_method = updatedPayment.Payment_method
	payment.Amount = updatedPayment.Amount
	payment.Payment_date = updatedPayment.Payment_date
	payment.Status = updatedPayment.Status
	if err := db.Where("id = ?", id).Updates(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func DeletePayment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&models.Payment{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	if err := db.Where("id = ?", id).Delete(&models.Payment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
