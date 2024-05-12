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

// GetPayment fetches a single payment by its ID from the URL parameters.
// It validates payment data and returns the payment details or an error message if the payment is not found or the data is invalid.
func GetPayment(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var payment models.Payment

	if err := db.Where("id = ?", id).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

// GetPayments retrieves all payments from the database.
// It returns a list of payments or an error message if the retrieval fails.
func GetPayments(c *gin.Context, db *gorm.DB) {
	payments, err := models.GetAllPayments(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving payments"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// SearchAllPayments retrieves all payments from the database based on the search parameters provided in the query string.
// It responds with a list of payments if successful or an informational message if no payments exist.
// On failure, it returns an HTTP 500 Internal Server Error.
// The search parameters include order_id, payment_method, amount, payment_date, and status.
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

// CreatePayment adds a new payment record to the database based on the JSON data provided in the request body.
// It validates the input data and responds with the created payment or an error message if the data is invalid or creation fails.
func CreatePayment(c *gin.Context, db *gorm.DB) {
	var newPayment models.Payment
	if err := c.ShouldBindJSON(&newPayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
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

// UpdatePayment modifies an existing payment record based on the JSON input and the ID provided in the URL.
// It checks the validity of the input data and updates the payment in the database, responding
// with the updated payment or an error message.
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

// DeletePayment removes a payment record from the database based on its ID provided in the URL.
// It handles the deletion process and responds with HTTP 204 No Content on success or an
// error message if the payment is not found or deletion fails.
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

// checkPayment validates the input data for a payment and returns an error if the data is invalid.
// It checks the payment's order_id, payment_method, amount, payment_date, and status fields for correct formatting.
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
