package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// GetBrand fetches a single brand based on the ID provided in the URL.
// It returns the brand if found or appropriate error messages for missing ID or not found scenarios.
// In case of an error, it sends an HTTP 500 Internal Server Error.
func GetBrand(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var brand models.Brands

	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}
	c.JSON(http.StatusOK, brand)
}

// GetBrands retrieves all brands from the database.
// It sends an HTTP 200 OK response with a list of brands or a message if no brands exist.
// In case of an error, it sends an HTTP 500 Internal Server Error.
func GetBrands(c *gin.Context, db *gorm.DB) {
	brands, err := models.GetAllBrands(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

// SearchAllBrands retrieves all brands from the database based on the search parameters provided in the query string.
// It responds with a list of brands if successful or an informational message if no brands exist.
// On failure, it returns an HTTP 500 Internal Server Error.
func SearchAllBrands(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"name", "description"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			searchParams[field] = cleanValue
		}
	}

	brands, err := models.SearchBrand(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve brands", "details": err.Error()})
		return
	}
	if !tools.CheckString(brands.Name, 250) && !tools.CheckString(brands.Description, 1000) {
		c.JSON(http.StatusNotFound, gin.H{"error": "No brand found"})
		return
	}

	c.JSON(http.StatusOK, brands)
}

// CreateBrand adds a new brand to the database based on the JSON data provided in the request body.
// It responds with the newly created brand or an error message if the data is invalid or creation fails.
// The brand's name and description fields are validated for correct formatting.
// If the brand is successfully created, it sends an HTTP 201 Created response.
// In case of a validation error, it sends an HTTP 400 Bad Request response.
func CreateBrand(c *gin.Context, db *gorm.DB) {
	var newBrand models.Brands
	if err := c.ShouldBindJSON(&newBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	brand := models.Brands{
		Name:        newBrand.Name,
		Description: newBrand.Description,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkBrand(brand, newBrand); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create brand", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, brand)
}

// UpdateBrand modifies an existing brand based on the ID provided in the URL.
// It updates the brand's name and description with the provided data and responds accordingly.
// If the brand is not found, it sends an HTTP 404 Not Found response.
// If the update is successful, it sends an HTTP 200 OK response with the updated brand.
// If the update fails, it sends an HTTP 500 Internal Server Error.
func UpdateBrand(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.BrandExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}

	var updatedBrand models.Brands
	if err := c.ShouldBindJSON(&updatedBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var brand models.Brands
	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}

	brand.Name = updatedBrand.Name
	brand.Description = updatedBrand.Description

	if failed, err := checkBrand(brand, updatedBrand); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update brand", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, brand)
}

// DeleteBrand removes a brand from the database based on the ID provided in the URL.
// It responds with an HTTP 204 No Content on success or an error message if the brand is not found or if deletion fails.
func DeleteBrand(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.BrandExists(db, convertedId) {
		fmt.Println("Brands does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Brands not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.Brands{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting brands"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// checkBrand validates the input data for a brand and returns an error if the data is invalid.
// It checks the brand's name and description fields for correct formatting.
func checkBrand(brand models.Brands, newBrand models.Brands) (bool, error) {
	switch true {
	case !brand.SetName(newBrand.Name):
		return true, fmt.Errorf("name is wrong formatted")
	case !brand.SetDescription(newBrand.Description):
		return true, fmt.Errorf("description is wrong formatted")
	}
	return false, nil
}
