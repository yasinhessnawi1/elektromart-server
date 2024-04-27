package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// GetCategories retrieves all categories from the database.
// Responds with a list of categories if successful or an informational message if no categories exist.
// On failure, it returns an HTTP 500 Internal Server Error.
func GetCategories(c *gin.Context, db *gorm.DB) {
	categories, err := models.GetAllCategories(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving categories"})
		return
	}
	if len(categories) == 0 {
		c.JSON(http.StatusOK, gin.H{"info": "No categories found, please create one first"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetCategory fetches a single category based on its ID provided in the URL path.
// It checks for valid category data and returns an HTTP 200 OK with the category details or an error if not found or data is invalid.
func GetCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if !tools.CheckString(category.Name, 255) || !tools.CheckString(category.Description, 1000) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid category data"})
		return
	}
	c.JSON(http.StatusOK, category)

}

// CreateCategory handles the creation of a new category via JSON input.
// It validates input and responds with the created category object or an error message on failure.
func CreateCategory(c *gin.Context, db *gorm.DB) {
	var newCategory models.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !tools.CheckString(newCategory.Name, 255) || !tools.CheckString(newCategory.Description, 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	fmt.Println(newCategory)
	var category models.Category
	category.Model.ID = uint(tools.GenerateUUID())
	category.Name = newCategory.Name
	category.Description = newCategory.Description

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// UpdateCategory modifies an existing category based on its ID.
// It validates the input data and updates the category in the database, responding with the updated data or an error.
func UpdateCategory(c *gin.Context, db *gorm.DB) {
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !tools.CheckString(updatedCategory.Name, 255) || !tools.CheckString(updatedCategory.Description, 1000) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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

// DeleteCategory removes a category from the database based on its ID.
// It handles the deletion process and returns an HTTP 204 No Content on success or an error message if the category is not found or deletion fails.
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
