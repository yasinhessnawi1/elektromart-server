package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetBrands(c *gin.Context, db *gorm.DB) {
	brands, err := models.GetAllBrands(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving brands"})
		return
	}
	c.JSON(http.StatusOK, brands)
}

func CreateBrand(c *gin.Context, db *gorm.DB) {
	var newBrand models.BrandsDB
	if err := c.ShouldBindJSON(&newBrand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var brand models.Brands
	brand.Brand_ID = tools.GenerateUUID()
	brand.Name = newBrand.Name
	brand.Description = newBrand.Description
	if err := db.Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, brand)
}

func UpdateBrand(c *gin.Context, db *gorm.DB) {
	var updatedBrand models.BrandsDB
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
	if err := db.Where("id = ?", id).Updates(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedBrand)
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
