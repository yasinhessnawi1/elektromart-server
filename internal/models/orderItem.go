package models

import (
	"gorm.io/gorm"
)

type OrderItem struct {
	BaseModel
	OrderID   uint
	ProductID uint
	Quantity  int
	Subtotal  float64
}

func GetAllOrderItems(db *gorm.DB) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := db.Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}
