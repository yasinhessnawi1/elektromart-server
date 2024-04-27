package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetBrands retrieves all brands from the database.
// It sends an HTTP 200 OK response with a list of brands or a message if no brands exist.
// In case of an error, it sends an HTTP 500 Internal Server Error.
func GetBrands(c *gin.Context, db *gorm.DB) {
	brands, err := models.GetAllBrands(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving brands", "error_message": err.Error()})
		return
	} else if len(brands) == 0 {
		c.JSON(http.StatusOK, gin.H{"info": "No brands found, please create one first"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

// GetBrand fetches a single brand based on the ID provided in the URL.
// It returns the brand if found or appropriate error messages for missing ID or not found scenarios.
func GetBrand(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID not provided"})
		return
	}
	var brand models.Brands
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}
	c.JSON(http.StatusOK, brand)
}

// CreateBrand handles the creation of a new brand.
// It validates the input and stores the new brand in the database.
// Responds with the created brand or an error message.
func CreateBrand(c *gin.Context, db *gorm.DB) {
	var newBrand models.Brands
	if err := c.ShouldBindJSON(&newBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var brand models.Brands
	brand.Model.ID = uint(tools.GenerateUUID())

	brand.Name = newBrand.Name
	brand.Description = newBrand.Description
	if err := db.Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, brand)
}

// UpdateBrand modifies an existing brand based on the ID provided in the URL.
// It updates the brand's name and description with the provided data and responds accordingly.
func UpdateBrand(c *gin.Context, db *gorm.DB) {
	var updatedBrand models.Brands
	if err := c.ShouldBindJSON(&updatedBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var brand models.Brands
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}
	brand.Name = updatedBrand.Name
	brand.Description = updatedBrand.Description
	if err := db.Where("id = ?", id).Updates(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedBrand)
}

// DeleteBrand removes a brand from the database based on the ID provided in the URL.
// It responds with an HTTP 204 No Content on success or an error message if the brand is not found or if deletion fails.
func DeleteBrand(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&models.Brands{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}
	if err := db.Where("id = ?", id).Delete(&models.Brands{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
