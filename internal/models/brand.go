package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

type Brands struct {
	gorm.Model
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

func (b *Brands) SetName(name string) bool {
	if !tools.CheckString(name, 255) {
		return false
	} else {
		b.Name = name
		return true
	}
}

func (b *Brands) SetDescription(description string) bool {
	if !tools.CheckString(description, 1000) {
		return false
	} else {
		b.Description = description
		return true
	}
}

func BrandExists(db *gorm.DB, id uint32) bool {
	var brand Brands
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		return false
	}
	return true
}

func SearchBrand(db *gorm.DB, searchParams map[string]interface{}) ([]Brands, error) {
	var brands []Brands
	query := db.Model(&Brands{})

	for key, value := range searchParams {
		switch key {
		case "name", "description":
			// For string fields
			if strVal, ok := value.(string); ok {
				query = query.Where(key+" LIKE ?", "%"+strVal+"%")
			}
		}
	}

	if err := query.Find(&brands).Debug().Error; err != nil {
		return nil, err
	}
	return brands, nil
}
