package models

import (
	"E-Commerce_Website_Database/internal/tools"
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

func (c *Category) SetName(name string) bool {
	if !tools.CheckString(name, 255) {
		return false
	} else {
		c.Name = name
		return true
	}
}

func (c *Category) SetDescription(description string) bool {
	if !tools.CheckString(description, 1000) {
		return false
	} else {
		c.Description = description
		return true
	}
}

func CategoryExists(db *gorm.DB, id uint32) bool {
	var category Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		return false
	}
	return true
}
