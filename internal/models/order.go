package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// Order represents the order model for transactions.
// It includes fields like User_ID, Order_date, Total_amount, and Status, which are tagged for JSON serialization.
type Order struct {
	gorm.Model
	User_ID      uint32  `json:"user_id"`
	Order_date   string  `json:"order_date"`
	Total_amount float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

// GetAllOrders retrieves all orders from the database.
// It returns a slice of Order objects or an error if there is any issue during fetching.
func GetAllOrders(db *gorm.DB) ([]Order, error) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// SetUserID sets the user ID for an order after verifying the existence of the user in the database.
// It does not return any value but updates the User_ID of the order if the user exists.
func (o *Order) SetUserID(user_id uint32, db *gorm.DB) {
	if UserExists(db, user_id) {
		o.User_ID = user_id
	}
}

// SetOrderDate validates and sets the order date.
// It returns true if the date is valid according to predefined date checking logic, otherwise returns false.
func (o *Order) SetOrderDate(order_date string) bool {
	if !tools.CheckDate(order_date) {
		return false
	}
	o.Order_date = order_date
	return true
}

// SetTotalAmount validates and sets the total amount of an order.
// It returns true if the amount is valid, otherwise returns false.
// The order's total amount is updated if the check is successful.
func (o *Order) SetTotalAmount(total_amount float64) bool {
	if !tools.CheckFloat(total_amount) {
		return false
	}
	o.Total_amount = total_amount
	return true
}

// SetStatus validates and sets the status of an order.
// It returns true if the status is valid according to predefined status checking logic, otherwise returns false.
// The order's status is updated if the check is successful.
func (o *Order) SetStatus(status string) bool {
	if !tools.CheckStatus(status, 255) {
		return false
	}
	o.Status = status
	return true
}

// OrderExists checks if an order exists in the database by its ID.
// It returns true if the order is found, otherwise returns false.
func OrderExists(db *gorm.DB, id uint32) bool {
	var order Order
	if db.Where("id = ?", id).First(&order).Error != nil {
		return false
	}
	return true
}
