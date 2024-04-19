package models

import (
	"gorm.io/gorm"
)

type Product struct {
	BaseModel
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	BrandID       uint
	CategoryID    uint
}

func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
