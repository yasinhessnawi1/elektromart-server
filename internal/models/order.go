package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

type OrderDB struct {
	gorm.Model
	ID           uint32  `json:"order_id"`
	User_ID      uint32  `json:"user_id"`
	Order_date   string  `json:"order_date"`
	Total_amount float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

type Order struct {
	gorm.Model
	User_ID      uint32  `json:"user_id"`
	Order_date   string  `json:"order_date"`
	Total_amount float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

func GetAllOrders(db *gorm.DB) ([]Order, error) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *Order) SetUserID(user_id uint32, db *gorm.DB) {
	if UserExists(db, user_id) {
		o.User_ID = user_id
	}
}

func (o *Order) SetOrderDate(order_date string) bool {
	if !tools.CheckDate(order_date) {
		return false
	}
	return true
}

func (o *Order) SetTotalAmount(total_amount float64) bool {
	if !tools.CheckFloat(total_amount) {
		return false
	}
	o.Total_amount = total_amount
	return true
}

func (o *Order) SetStatus(status string) bool {
	if !tools.CheckStatus(status, 255) {
		return false
	}
	o.Status = status
	return true
}

func OrderExists(db *gorm.DB, id uint32) bool {
	var order Order
	if db.Where("id = ?", id).First(&order).Error != nil {
		return false
	}
	return true
}
