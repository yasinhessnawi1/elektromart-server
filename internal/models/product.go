package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// Product represents the product entity with properties such as name, description,
// price, stock quantity, and associations with brand and category.
type Product struct {
	gorm.Model
	Name           string  `json:"name"`           // Name of the product
	Description    string  `json:"description"`    // Description of the product
	Price          float64 `json:"price"`          // Price of the product
	Stock_quantity int     `json:"stock_quantity"` // Quantity in stock
	Brand_ID       uint32  `json:"brand_id"`       // Foreign key for the associated brand
	Category_ID    uint32  `json:"category_id"`    // Foreign key for the associated category
}

// GetAllProducts retrieves all products from the database.
// It returns a slice of Product or an error if the fetch fails.
func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// SetName sets the name of the product after validating its length.
// Returns true if the name is within the allowed length, otherwise false.
func (p *Product) SetName(name string) bool {
	if !tools.CheckString(name, 255) {
		return false
	} else {
		p.Name = name
		return true
	}
}

// SetDescription sets the product's description after validating its length.
// Returns true if the description is within the allowed length, otherwise false.
func (p *Product) SetDescription(description string) bool {
	if !tools.CheckString(description, 1000) {
		return false
	} else {
		p.Description = description
		return true
	}
}

// SetPrice sets the price of the product after validating it as a positive float.
// Returns true if the price is valid, otherwise false.
func (p *Product) SetPrice(price float64) bool {
	if !tools.CheckFloat(price) {
		return false
	} else {
		p.Price = price
		return true
	}
}

// SetStockQuantity sets the stock quantity of the product after validating it as a non-negative integer.
// Returns true if the stock quantity is valid, otherwise false.
func (p *Product) SetStockQuantity(stock_quantity int) bool {
	if !tools.CheckInt(stock_quantity) {
		return false
	} else {
		p.Stock_quantity = stock_quantity
		return true
	}
}

// SetBrandID sets the brand ID of the product, verifying the existence of the brand.
// Returns true if the brand exists, otherwise false.
func (p *Product) SetBrandID(brand_id uint32, db *gorm.DB) bool {
	if !BrandExists(db, brand_id) {
		return false
	} else {
		p.Brand_ID = brand_id
		return true
	}
}

// SetCategoryID sets the category ID of the product, verifying the existence of the category.
// Returns true if the category exists, otherwise false.
func (p *Product) SetCategoryID(category_id uint32, db *gorm.DB) bool {
	if !CategoryExists(db, category_id) {
		return false
	} else {
		p.Category_ID = category_id
		return true
	}
}

// ProductExists checks if a specific product exists in the database by its ID.
// Returns true if the product exists, otherwise false.
func ProductExists(db *gorm.DB, id uint32) bool {
	var product Product
	if db.First(&product, id).Error != nil {
		return false
	}
	return true
}

// SearchProduct performs a search based on given search parameters.
// It returns a slice of products that match the criteria or an error if the search fails.
func SearchProduct(db *gorm.DB, searchParams map[string]interface{}) ([]Product, error) {
	var products []Product
	query := db.Model(&Product{})

	for key, value := range searchParams {
		switch key {
		case "name", "description":
			if strVal, ok := value.(string); ok {
				query = query.Where(key+" LIKE ?", "%"+strVal+"%")
			}
		case "price", "stock_quantity", "brand_id", "category_id":
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
