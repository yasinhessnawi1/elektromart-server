package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetPayments retrieves all payments from the database.
// It returns a list of payments or an error message if the retrieval fails.
func GetPayments(c *gin.Context, db *gorm.DB) {
	payments, err := models.GetAllPayments(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving payments", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// GetPayment fetches a single payment by its ID from the URL parameters.
// It validates payment data and returns the payment details or an error message if the payment is not found or the data is invalid.
func GetPayment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var payment models.Payment
	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	if !tools.CheckInt(int(payment.Order_ID)) || !tools.CheckString(payment.Payment_method, 255) || !tools.CheckFloat(payment.Amount) || !tools.CheckDate(payment.Payment_date) || !tools.CheckString(payment.Status, 1000) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid payment data"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

// CreatePayment handles the creation of a new payment record from JSON input.
// It validates the input and stores the new payment in the database, responding with the created payment or an error message.
func CreatePayment(c *gin.Context, db *gorm.DB) {
	var newPayment models.Payment
	if err := c.ShouldBindJSON(&newPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !tools.CheckInt(int(newPayment.Order_ID)) || !tools.CheckString(newPayment.Payment_method, 255) || !tools.CheckFloat(newPayment.Amount) || !tools.CheckDate(newPayment.Payment_date) || !tools.CheckString(newPayment.Status, 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var payment models.Payment
	payment.Model.ID = uint(tools.GenerateUUID())
	payment.Order_ID = newPayment.Order_ID
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

// UpdatePayment modifies an existing payment record based on the JSON input and the ID provided in the URL.
// It checks the validity of the input data and updates the payment in the database, responding with the updated payment or an error message.
func UpdatePayment(c *gin.Context, db *gorm.DB) {
	var updatedPayment models.Payment
	if err := c.ShouldBindJSON(&updatedPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !tools.CheckInt(int(updatedPayment.Order_ID)) || !tools.CheckString(updatedPayment.Payment_method, 255) || !tools.CheckFloat(updatedPayment.Amount) || !tools.CheckDate(updatedPayment.Payment_date) || !tools.CheckString(updatedPayment.Status, 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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

// DeletePayment removes a payment record from the database based on its ID provided in the URL.
// It handles the deletion process and responds with HTTP 204 No Content on success or an error message if the payment is not found or deletion fails.
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
