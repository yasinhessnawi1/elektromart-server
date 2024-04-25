package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

type BrandsDB struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Brands struct {
	gorm.Model
	ID          uint32 `json:"id"`
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
