package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
	"strings"
)

// Payment represents the payment model associated with an order.
// It includes fields for the order ID, payment method, amount, payment date, and status, all of which include JSON serialization tags.
type Payment struct {
	gorm.Model
	Order_ID       uint32  `json:"order_id"`
	Payment_method string  `json:"payment_method"`
	Amount         float64 `json:"amount"`
	Payment_date   string  `json:"payment_date"`
	Status         string  `json:"status"`
}

// GetAllPayments retrieves all payments from the database.
// It returns a slice of Payment objects and an error if there is any issue during fetching.
func GetAllPayments(db *gorm.DB) ([]Payment, error) {
	var payments []Payment
	if err := db.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// SetOrderID sets the order ID for the payment after verifying the existence of the order.
// Returns true if the order exists and the ID is set; otherwise, it returns false.
func (p *Payment) SetOrderID(order_id uint32, db *gorm.DB) bool {
	if !OrderExists(db, order_id) {
		return false
	} else {
		p.Order_ID = order_id
		return true
	}
}

// SetPaymentMethod sets the payment method for the payment.
// It validates the payment method to ensure it meets certain criteria (not implemented here) and returns true if valid.
func (p *Payment) SetPaymentMethod(payment_method string) bool {
	if !tools.CheckPaymentMethod(payment_method) {
		return false
	} else {
		p.Payment_method = payment_method
		return true
	}
}

// SetAmount sets the amount of the payment.
// It ensures the amount is valid as a float value and returns true if so; otherwise, it returns false.
func (p *Payment) SetAmount(amount float64) bool {
	if !tools.CheckFloat(amount) {
		return false
	} else {
		p.Amount = amount
		return true
	}
}

// SetPaymentDate sets the date of the payment.
// It validates the date format and returns true if the date is valid; otherwise, it returns false.
func (p *Payment) SetPaymentDate(payment_date string) bool {
	if !tools.CheckDate(payment_date) {
		return false
	} else {
		p.Payment_date = payment_date
		return true
	}
}

// SetStatus sets the status of the payment.
// It validates the status based on predefined criteria and returns true if the status is valid; otherwise, it returns false.
func (p *Payment) SetStatus(status string) bool {
	if !tools.CheckStatus(status, 255) {
		return false
	} else {
		p.Status = status
		return true
	}
}

// PaymentExists checks if a payment exists in the database by its ID.
// It returns true if the payment is found, otherwise returns false.
func PaymentExists(db *gorm.DB, id uint32) bool {
	var payment Payment
	if db.Where("id = ?", id).First(&payment).Error != nil {
		return false
	}
	return true
}

// SearchPayment performs a search for a payment based on the provided search parameters.
// It constructs a search query dynamically and returns the matching payment or an error if not found.
// If the search is successful, it returns the payment.
// If no payment is found, it responds with an HTTP 404 Not Found status.
// If the search is successful, it responds with an HTTP 200 OK status and the payment details in JSON format.
func SearchPayment(db *gorm.DB, searchParams map[string]interface{}) ([]Payment, error) {
	var payments []Payment
	query := db.Model(&Payment{})

	for key, value := range searchParams {
		valueStr, isString := value.(string)
		switch key {
		case "payment_method", "status":
			// For string fields
			if isString {
				query = query.Where(key+" LIKE ?", "%"+strings.ToLower(valueStr)+"%")
			}
		case "amount":
			// For numeric fields
			if numVal, ok := value.(float64); ok {
				query = query.Where(key+" = ?", numVal)
			}
		case "order_id":
			if numVal, ok := value.(int); ok {
				query = query.Where(key+" = ?", numVal)
			}

		case "payment_date":
			if isString {
				query = query.Where(key+" = ?", valueStr)
			}
		}
	}

	if err := query.Find(&payments).Debug().Error; err != nil {
		return nil, err
	}
	return payments, nil
}
