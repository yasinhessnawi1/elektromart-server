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

func GetUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	fmt.Println(id)
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Limit(1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// in /internal/handlers/user_handler.go
func CreateUser(c *gin.Context, db *gorm.DB) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	user.Model.ID = uint(tools.GenerateUUID())
	if checkUser(c, user, newUser) {
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func checkUser(c *gin.Context, user models.User, newUser models.User) bool {
	switch true {
	case !user.SetFirstName(newUser.First_Name):
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name is wrong formatted"})
		return true
	case !user.SetLastName(newUser.Last_Name):
		c.JSON(http.StatusBadRequest, gin.H{"error": "last name is wrong formatted"})
		return true
	case !user.SetUsername(newUser.Username):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
		return true
	case !user.SetPassword(newUser.Password):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return true
	case !user.SetEmail(newUser.Email):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return true
	case !user.SetAddress(newUser.Address):
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address"})
		return true
	}
	return false
}
func UpdateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	id := c.Param("id")
	if !models.UserExists(db, tools.ConvertStringToUint(id)) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if checkUser(c, user, newUser) {
		return
	}
	db.Save(&user)
	c.JSON(http.StatusOK, user)
}
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if !models.UserExists(db, tools.ConvertStringToUint(id)) {
		return
	}
	if err := db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
