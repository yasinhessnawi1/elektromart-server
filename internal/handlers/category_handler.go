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

func SearchAllCategories(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"name", "description"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			searchParams[field] = cleanValue
		}
	}

	categories, err := models.SearchCategory(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories", "details": err.Error()})
		return
	}

	if len(categories) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No category found"})
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

	if failed, err := checkCategory(category, newCategory); failed {
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
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.CategoryExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.Name = updatedCategory.Name
	category.Description = updatedCategory.Description

	if failed, err := checkCategory(category, updatedCategory); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.CategoryExists(db, convertedId) {
		fmt.Println("Category does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.Category{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting category"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func checkCategory(category models.Category, newCategory models.Category) (bool, error) {
	switch true {
	case !category.SetName(newCategory.Name):
		return true, fmt.Errorf("name is wrong formatted")
	case !category.SetDescription(newCategory.Description):
		return true, fmt.Errorf("description is wrong formatted")
	}
	return false, nil
}
