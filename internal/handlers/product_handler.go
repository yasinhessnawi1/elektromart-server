package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetProducts(c *gin.Context, db *gorm.DB) {
	products, err := models.GetAllProducts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
func CreateProduct(c *gin.Context, db *gorm.DB) {
	var newProduct models.ProductDB
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product models.Product
	product.Product_ID = tools.GenerateUUID()
	product.Name = newProduct.Name
	product.Description = newProduct.Description
	product.Price = newProduct.Price
	product.Stock_quantity = newProduct.Stock_quantity
	product.Brand_ID = newProduct.Brand_ID
	product.Category_ID = newProduct.Category_ID

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newProduct)
}
func UpdateProduct(c *gin.Context, db *gorm.DB) {
	var product models.Product
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	var newProduct models.ProductDB
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.Name = newProduct.Name
	product.Description = newProduct.Description
	product.Price = newProduct.Price
	product.Stock_quantity = newProduct.Stock_quantity
	product.Brand_ID = newProduct.Brand_ID
	product.Category_ID = newProduct.Category_ID
	db.Save(&product)
	c.JSON(http.StatusOK, product)
}
func DeleteProduct(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
