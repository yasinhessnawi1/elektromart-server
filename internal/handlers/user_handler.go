package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetUsers(c *gin.Context, db *gorm.DB) {
	users, err := models.GetAllUsers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// in /internal/handlers/user_handler.go
func CreateUser(c *gin.Context, db *gorm.DB) {
	var newUser models.UserDB
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	user.User_ID = tools.GenerateUUID()
	user.Username = newUser.Username
	user.Password = newUser.Password
	user.Email = newUser.Email
	user.First_Name = newUser.First_Name
	user.Last_Name = newUser.Last_Name
	user.Address = newUser.Address
	fmt.Println(user)
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
func UpdateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var newUser models.UserDB
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Username = newUser.Username
	user.Password = newUser.Password
	user.Email = newUser.Email
	user.First_Name = newUser.First_Name
	user.Last_Name = newUser.Last_Name
	user.Address = newUser.Address
	db.Save(&user)
	c.JSON(http.StatusOK, user)
}
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
