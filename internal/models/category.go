package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// Category represents the category model for products.
// It extends gorm.Model, adding Name and Description fields with JSON tags to aid in serialization.
type Category struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetAllCategories retrieves all categories from the database.
// It returns a slice of Category and an error if there is any issue during fetching.
func GetAllCategories(db *gorm.DB) ([]Category, error) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// SetName attempts to set the category's name while enforcing a maximum length of 255 characters.
// It returns true if the name is valid and set successfully, otherwise returns false.
func (c *Category) SetName(name string) bool {
	if !tools.CheckString(name, 255) {
		return false
	} else {
		c.Name = name
		return true
	}
}

// SetDescription attempts to set the category's description while enforcing a maximum length of 1000 characters.
// It returns true if the description is valid and set successfully, otherwise returns false.
func (c *Category) SetDescription(description string) bool {
	if !tools.CheckString(description, 1000) {
		return false
	} else {
		c.Description = description
		return true
	}
}

// CategoryExists checks the existence of a category by its ID in the database.
// It returns true if the category is found, otherwise false if the category does not exist or there is an error.
func CategoryExists(db *gorm.DB, id uint32) bool {
	var category Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		return false
	}
	return true
}

// SearchCategory performs a search for a category based on the provided search parameters.
// It constructs a search query dynamically and returns the matching category or an error if not found.
func SearchCategory(db *gorm.DB, searchParams map[string]interface{}) (Category, error) {
	var category Category
	query := db.Model(&Category{})

	for key, value := range searchParams {
		switch key {
		case "name", "description":
			// For string fields
			if strVal, ok := value.(string); ok {
				query = query.Where(key+" = ?", strVal)
			}
		}
	}

	if err := query.First(&category).Debug().Error; err != nil {
		return category, err
	}
	return category, nil
}
