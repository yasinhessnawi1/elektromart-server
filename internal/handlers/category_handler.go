package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

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
