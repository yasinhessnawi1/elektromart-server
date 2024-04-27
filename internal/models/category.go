package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
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

func SearchCategory(db *gorm.DB, searchParams map[string]interface{}) ([]Category, error) {
	var categories []Category
	query := db.Model(&Category{})

	for key, value := range searchParams {
		switch key {
		case "name", "description":
			// For string fields
			if strVal, ok := value.(string); ok {
				query = query.Where(key+" LIKE ?", "%"+strVal+"%")
			}
		}
	}

	if err := query.Find(&categories).Debug().Error; err != nil {
		return nil, err
	}
	return categories, nil
}
