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
// If the fetch is successful, it returns the list of brands, otherwise an error message.
func GetAllBrands(db *gorm.DB) ([]Brands, error) {
	var brands []Brands
	if err := db.Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}

// SetName sets the name of the brand with validation.
// It ensures the name does not exceed 255 characters. Returns true if set successfully, otherwise false.
// It returns true if the name is within the allowed length, otherwise false.
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
// It returns true if the description is within the allowed length, otherwise false.
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
// It returns true if the brand exists, otherwise false.
func BrandExists(db *gorm.DB, id uint32) bool {
	var brand Brands
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		return false
	}
	return true
}

// SearchBrand performs a search on brands based on provided query parameters.
// It constructs a search query dynamically and returns the matching brand or an appropriate error message.
// If no brand is found, it responds with an HTTP 404 Not Found status.
// If the search is successful, it responds with an HTTP 200 OK status and the brand details in JSON format.
func SearchBrand(db *gorm.DB, searchParams map[string]interface{}) (Brands, error) {
	var brands Brands
	query := db.Model(&Brands{})

	for key, value := range searchParams {
		switch key {
		case "name", "description":
			// For string fields
			if strVal, ok := value.(string); ok {
				query = query.Where(key+" = ?", strVal)
			}
		}
	}

	if err := query.First(&brands).Debug().Error; err != nil {
		return brands, err
	}
	return brands, nil
}
