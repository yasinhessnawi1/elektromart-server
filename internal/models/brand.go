package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// Brands represents the brand model that holds details about a brand.
// It includes the default gorm.Model fields along with Name and Description for the brand.
type Brands struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetAllBrands retrieves all brands from the database.
// It returns a slice of Brands and an error if there is any issue in fetching the data.
func GetAllBrands(db *gorm.DB) ([]Brands, error) {
	var brands []Brands
	if err := db.Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}

// SetName sets the name of the brand with validation.
// It ensures the name does not exceed 255 characters. Returns true if set successfully, otherwise false.
func (b *Brands) SetName(name string) bool {
	if !tools.CheckString(name, 255) {
		return false
	} else {
		b.Name = name
		return true
	}
}

// SetDescription sets the description of the brand with validation.
// It ensures the description does not exceed 1000 characters. Returns true if set successfully, otherwise false.
func (b *Brands) SetDescription(description string) bool {
	if !tools.CheckString(description, 1000) {
		return false
	} else {
		b.Description = description
		return true
	}
}

// BrandExists checks if a brand exists in the database by its ID.
// It queries the database for the brand by the given ID and returns true if found, otherwise false.
func BrandExists(db *gorm.DB, id uint32) bool {
	var brand Brands
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		return false
	}
	return true
}
