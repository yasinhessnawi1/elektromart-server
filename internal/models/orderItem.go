package models

import (
	"gorm.io/gorm"
)

type OrderItemDB struct {
	gorm.Model
	Order_ID   uint32  `json:"order_id"`
	Product_ID uint32  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Subtotal   float64 `json:"subtotal"`
}

type OrderItem struct {
	gorm.Model
	ID         uint32  `json:"id"`
	Order_ID   uint32  `json:"order_id"`
	Product_ID uint32  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Subtotal   float64 `json:"subtotal"`
}

func GetAllOrderItems(db *gorm.DB) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := db.Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}
