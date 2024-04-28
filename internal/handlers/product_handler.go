package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func GetProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var product models.Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func GetProducts(c *gin.Context, db *gorm.DB) {
	products, err := models.GetAllProducts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving products", "error_message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func SearchAllProducts(c *gin.Context, db *gorm.DB) {
	queryParam := c.Param("any") // Get the query parameter named 'any'

	// Check if the query parameter is empty
	if queryParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'any' is required"})
		return
	}

	queryParam = strings.TrimSpace(queryParam) // Clean the input
	searchParams := make(map[string]interface{})

	// Define the fields to search in
	searchFields := []string{"name", "description", "price", "stock_quantity", "brand_id", "category_id"}

	// Populate the searchParams map with the same queryParam for all specified fields
	for _, field := range searchFields {
		switch field {
		case "price", "stock_quantity", "brand_id", "category_id":
			if numVal, err := strconv.ParseFloat(queryParam, 64); err == nil {
				searchParams[field] = numVal // Use as a number if conversion is successful
			} // If conversion fails, do not add to map - optional: handle this as needed
		default:
			searchParams[field] = queryParam // Use as a string
		}
	}

	// Search for products with the populated parameters
	products, err := models.SearchProduct(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products", "details": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
		return
	}

	c.JSON(http.StatusOK, products)
}
func CreateProduct(c *gin.Context, db *gorm.DB) {
	var newProduct models.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:           newProduct.Name,
		Description:    newProduct.Description,
		Price:          newProduct.Price,
		Stock_quantity: newProduct.Stock_quantity,
		Brand_ID:       newProduct.Brand_ID,
		Category_ID:    newProduct.Category_ID,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkProduct(product, newProduct, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}
func UpdateProduct(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.ProductExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var newProduct models.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var product models.Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found during update"})
		return
	}
	product.Name = newProduct.Name
	product.Description = newProduct.Description
	product.Price = newProduct.Price
	product.Stock_quantity = newProduct.Stock_quantity
	product.Brand_ID = newProduct.Brand_ID
	product.Category_ID = newProduct.Category_ID

	if failed, err := checkProduct(product, newProduct, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
func DeleteProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.ProductExists(db, convertedId) {
		fmt.Println("Product does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.Product{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Product deleted"})
}

func checkProduct(product models.Product, newProduct models.Product, db *gorm.DB) (bool, error) {
	switch true {
	case !product.SetName(newProduct.Name):
		return true, fmt.Errorf("name is wrong formatted")
	case !product.SetDescription(newProduct.Description):
		return true, fmt.Errorf("description is wrong formatted")
	case !product.SetPrice(newProduct.Price):
		return true, fmt.Errorf("invalid price")
	case !product.SetStockQuantity(newProduct.Stock_quantity):
		return true, fmt.Errorf("invalid stock quantity")
	case !product.SetBrandID(newProduct.Brand_ID, db):
		return true, fmt.Errorf("invalid brand_id or not existing")
	case !product.SetCategoryID(newProduct.Category_ID, db):
		return true, fmt.Errorf("invalid category_id or not existing")
	}
	return false, nil
}
