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

// GetProduct retrieves a single product by its ID.
// It checks for the product's existence and validity of its data, then returns the product details or an error message.
func GetProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var product models.Product

	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetProducts retrieves all products from the database.
// It returns a JSON response with a list of products or an error message if the retrieval fails.
func GetProducts(c *gin.Context, db *gorm.DB) {
	products, err := models.GetAllProducts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// SearchAllProducts performs a search on products based on provided query parameters.
// It constructs a search query dynamically and returns the matching products or an appropriate error message.
func SearchAllProducts(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"name", "description", "price", "stock_quantity", "brand_name", "category_name"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			if field == "price" || field == "stock_quantity" {
				if numVal, err := strconv.ParseFloat(cleanValue, 64); err == nil {
					searchParams[field] = numVal
				}
			} else {
				searchParams[field] = cleanValue
			}
		}
	}

	// Search for products based on the search parameters
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

// CreateProduct handles the creation of a new product from JSON input.
// It validates the input and stores the new product in the database, responding with the created product or an error message.
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

// UpdateProduct handles the updating of an existing product.
// It validates the provided input and updates the product in the database, responding with the updated product or an error message.
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

// DeleteProduct handles the deletion of a product by its ID.
// It validates the product's existence and removes it from the database, responding with an appropriate message.
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

// checkProduct performs validation checks on product data.
// It returns a boolean indicating failure and an error with the validation issue.
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
