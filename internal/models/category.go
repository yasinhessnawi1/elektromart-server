package models

import (
	"gorm.io/gorm"
)

type CategoryDB struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Category struct {
	gorm.Model
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetAllCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
