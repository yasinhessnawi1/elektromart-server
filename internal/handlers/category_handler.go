package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var category models.Category

	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)

}

func GetCategories(c *gin.Context, db *gorm.DB) {
	categories, err := models.GetAllCategories(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context, db *gorm.DB) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:        newCategory.Name,
		Description: newCategory.Description,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkCategory(category, newCategory, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context, db *gorm.DB) {
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var category models.Category
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	category.Name = updatedCategory.Name
	category.Description = updatedCategory.Description
	if err := db.Where("id = ?", id).Updates(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedCategory)
}

func DeleteCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&models.Category{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err := db.Where("id = ?", id).Delete(&models.Category{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func checkCategory(category models.Category, newCategory models.Category, db *gorm.DB) (bool, error) {
	switch true {
	case !category.SetName(newCategory.Name):
		return true, fmt.Errorf("name is wrong formatted")
	case !category.SetDescription(newCategory.Description):
		return true, fmt.Errorf("description is wrong formatted")
	}
	return false, nil
}
