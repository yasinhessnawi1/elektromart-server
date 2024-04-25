package models

import (
	"gorm.io/gorm"
)

type ProductDB struct {
	gorm.Model
	Name           string  ` json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Stock_quantity int     `json:"stock_quantity"`
	Brand_ID       uint32  `json:"brand_id"`
	Category_ID    uint32  `json:"category_id"`
}

type Product struct {
	gorm.Model
	Product_ID     uint32  `json:"product_id"`
	Name           string  ` json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Stock_quantity int     `json:"stock_quantity"`
	Brand_ID       uint32  `json:"brand_id"`
	Category_ID    uint32  `json:"category_id"`
}

func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
