package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetBrand(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var brand models.Brands

	if err := db.Where("id = ?", id).First(&brand).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		return
	}
	c.JSON(http.StatusOK, brand)
}

func GetBrands(c *gin.Context, db *gorm.DB) {
	brands, err := models.GetAllBrands(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

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

func checkBrand(brand models.Brands, newBrand models.Brands) (bool, error) {
	switch true {
	case !brand.SetName(newBrand.Name):
		return true, fmt.Errorf("name is wrong formatted")
	case !brand.SetDescription(newBrand.Description):
		return true, fmt.Errorf("description is wrong formatted")
	}
	return false, nil
}
