package models

import (
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
