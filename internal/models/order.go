package models

import (
	"gorm.io/gorm"
)

type Order struct {
	BaseModel
	UserID      uint
	OrderDate   string
	TotalAmount float64
	Status      string
}

func GetAllOrders(db *gorm.DB) ([]Order, error) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
