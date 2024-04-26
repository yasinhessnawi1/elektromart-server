package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

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

func (oi *OrderItem) SetOrderID(order_id uint32, db *gorm.DB) bool {
	if !OrderExists(db, order_id) {
		return false
	} else {
		oi.Order_ID = order_id
		return true
	}
}

func (oi *OrderItem) SetProductID(product_id uint32, db *gorm.DB) bool {
	if !ProductExists(db, product_id) {
		return false
	} else {
		oi.Product_ID = product_id
		return true
	}
}

func (oi *OrderItem) SetQuantity(quantity int) bool {
	if !tools.CheckInt(quantity) {
		return false
	} else {
		oi.Quantity = quantity
		return true
	}
}

func (oi *OrderItem) SetSubtotal(subtotal float64) bool {
	if !tools.CheckFloat(subtotal) {
		return false
	} else {
		oi.Subtotal = subtotal
		return true
	}
}

func OrderItemExists(db *gorm.DB, id uint32) bool {
	var orderItem OrderItem
	if db.Where("id = ?", id).First(&orderItem).Error != nil {
		return false
	}
	return true
}
