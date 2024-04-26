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

func GetPayment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var payment models.Payment

	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func GetPayments(c *gin.Context, db *gorm.DB) {
	payments, err := models.GetAllPayments(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving payments"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func SearchAllPayments(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"order_id", "payment_method", "amount", "payment_date", "status"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			switch field {
			case "order_id":
				if numVal, err := strconv.Atoi(cleanValue); err == nil {
					searchParams[field] = numVal
				}
			case "amount":
				if numVal, err := strconv.ParseFloat(cleanValue, 64); err == nil {
					searchParams[field] = numVal
				}
			case "payment_method", "status":
				searchParams[field] = strings.ToLower(cleanValue)
			case "payment_date":
				searchParams[field] = cleanValue
			default:
				searchParams[field] = cleanValue
			}
		}
	}

	payments, err := models.SearchPayment(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment", "details": err.Error()})
		return
	}

	if len(payments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No payments found"})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func CreatePayment(c *gin.Context, db *gorm.DB) {
	var newPayment models.Payment
	if err := c.ShouldBindJSON(&newPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := models.Payment{
		Order_ID:       newPayment.Order_ID,
		Payment_method: newPayment.Payment_method,
		Amount:         newPayment.Amount,
		Payment_date:   newPayment.Payment_date,
		Status:         newPayment.Status,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkPayment(payment, newPayment, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func UpdatePayment(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.PaymentExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	var updatedPayment models.Payment
	if err := c.ShouldBindJSON(&updatedPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var payment models.Payment
	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	payment.Order_ID = updatedPayment.Order_ID
	payment.Payment_method = updatedPayment.Payment_method
	payment.Amount = updatedPayment.Amount
	payment.Payment_date = updatedPayment.Payment_date
	payment.Status = updatedPayment.Status

	if failed, err := checkPayment(payment, updatedPayment, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func DeletePayment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.PaymentExists(db, convertedId) {
		fmt.Println("Payment does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.Payment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting payment"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func checkPayment(payment models.Payment, newPayment models.Payment, db *gorm.DB) (bool, error) {
	switch true {
	case !payment.SetOrderID(newPayment.Order_ID, db):
		return true, fmt.Errorf("invalid order_id or not existing")
	case !payment.SetPaymentMethod(newPayment.Payment_method):
		return true, fmt.Errorf("payment metode is not expected")
	case !payment.SetAmount(newPayment.Amount):
		return true, fmt.Errorf("invalid amount")
	case !payment.SetPaymentDate(newPayment.Payment_date):
		return true, fmt.Errorf("invalid payment date")
	case !payment.SetStatus(newPayment.Status):
		return true, fmt.Errorf("payment status is not expected")
	}
	return false, nil
}
