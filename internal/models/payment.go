package models

import (
	"gorm.io/gorm"
)

type PaymentDB struct {
	gorm.Model
	Order_ID       uint32  `json:"order_id"`
	Payment_method string  `json:"payment_method"`
	Amount         float64 `json:"amount"`
	Payment_date   string  `json:"payment_date"`
	Status         string  `json:"status"`
}

type Payment struct {
	gorm.Model
	Payment_ID     uint32  `json:"payment_id"`
	Order_ID       uint32  `json:"order_id"`
	Payment_method string  `json:"payment_method"`
	Amount         float64 `json:"amount"`
	Payment_date   string  `json:"payment_date"`
	Status         string  `json:"status"`
}

func GetAllPayments(db *gorm.DB) ([]Payment, error) {
	var payments []Payment
	if err := db.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
