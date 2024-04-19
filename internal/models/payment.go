package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	BaseModel
	OrderID       uint
	PaymentMethod string
	Amount        float64
	PaymentDate   string
	Status        string
}

func GetAllPayments(db *gorm.DB) ([]Payment, error) {
	var payments []Payment
	if err := db.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
