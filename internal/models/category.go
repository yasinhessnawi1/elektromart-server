package models

import (
	"gorm.io/gorm"
)

type Category struct {
	BaseModel
	Name        string
	Description string
}

func GetAllCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
