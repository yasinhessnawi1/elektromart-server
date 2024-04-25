package models

import (
	"gorm.io/gorm"
)

type BrandsDB struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Brands struct {
	gorm.Model
	Brand_ID    uint32 `json:"brand_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetAllBrands(db *gorm.DB) ([]Brands, error) {
	var brands []Brands
	if err := db.Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}
