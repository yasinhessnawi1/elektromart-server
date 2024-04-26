package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
	"strings"
)

type Payment struct {
	gorm.Model
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

func (p *Payment) SetOrderID(order_id uint32, db *gorm.DB) bool {
	if !OrderExists(db, order_id) {
		return false
	} else {
		p.Order_ID = order_id
		return true
	}
}

func (p *Payment) SetPaymentMethod(payment_method string) bool {
	if !tools.CheckPaymentMethod(payment_method) {
		return false
	} else {
		p.Payment_method = payment_method
		return true
	}
}

func (p *Payment) SetAmount(amount float64) bool {
	if !tools.CheckFloat(amount) {
		return false
	} else {
		p.Amount = amount
		return true
	}
}

func (p *Payment) SetPaymentDate(payment_date string) bool {
	if !tools.CheckDate(payment_date) {
		return false
	} else {
		p.Payment_date = payment_date
		return true
	}
}

func (p *Payment) SetStatus(status string) bool {
	if !tools.CheckStatus(status, 255) {
		return false
	} else {
		p.Status = status
		return true
	}
}

func PaymentExists(db *gorm.DB, id uint32) bool {
	var payment Payment
	if db.Where("id = ?", id).First(&payment).Error != nil {
		return false
	}
	return true
}

func SearchPayment(db *gorm.DB, searchParams map[string]interface{}) ([]Payment, error) {
	var products []Payment
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

	if err := query.Find(&products).Debug().Error; err != nil {
		return nil, err
	}
	return products, nil
}
