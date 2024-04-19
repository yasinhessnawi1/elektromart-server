package models

import (
	"gorm.io/gorm"
)

type Brand struct {
	BaseModel
	Name        string
	Description string
}

func GetAllBrands(db *gorm.DB) ([]Brand, error) {
	var brands []Brand
	if err := db.Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}
