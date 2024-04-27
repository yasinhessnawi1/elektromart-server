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

// GetUser retrieves a single user by ID from the URL parameters.
// It returns the user details or an error message if the user is not found.
func GetUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Limit(1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsers retrieves all users from the database.
// It returns a list of users or an error message if the retrieval fails.
func GetUsers(c *gin.Context, db *gorm.DB) {
	users, err := models.GetAllUsers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// SearchAllUsers performs a search for users based on provided query parameters.
// It constructs a search query dynamically and returns the matching users or an appropriate error message.
func SearchAllUsers(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}
	for _, field := range []string{"username", "email", "first_name", "last_name", "address"} {
		if value := c.Query(field); value != "" {
			searchParams[field] = strings.TrimSpace(value)
		}
	}
	users, err := models.SearchUsers(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser handles the creation of a new user from JSON input.
// It validates the input and stores the new user in the database, responding with the created user or an error message.
func CreateUser(c *gin.Context, db *gorm.DB) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}
	if !tools.CheckString(newUser.Username, 255) || !tools.CheckString(newUser.Password, 255) ||
		!tools.CheckString(newUser.Email, 255) || !tools.CheckString(newUser.First_Name, 255) ||
		!tools.CheckString(newUser.Last_Name, 255) || !tools.CheckString(newUser.Address, 255) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user := models.User{
		Username:   newUser.Username,
		Password:   newUser.Password,
		Email:      newUser.Email,
		First_Name: newUser.First_Name,
		Last_Name:  newUser.Last_Name,
		Address:    newUser.Address,
		Model:      gorm.Model{ID: uint(tools.GenerateUUID())},
	}
	if failed, err := checkUser(user, newUser); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// UpdateUser handles updating an existing user.
// It validates the user's existence and the provided input, then updates the user in the database.
func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))
	if !models.UserExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found during update"})
		return
	}
	user.Username = newUser.Username
	user.Password = newUser.Password
	user.Email = newUser.Email
	user.First_Name = newUser.First_Name
	user.Last_Name = newUser.Last_Name
	user.Address = newUser.Address
	if failed, err := checkUser(user, newUser); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles the deletion of a user by ID.
// It validates the user's existence and removes the user from the database, responding with an appropriate message.
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)
	if !models.UserExists(db, convertedId) {
		fmt.Println("User does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted"})
}

// checkUser performs validation checks on user data.
// It returns a boolean indicating failure and an error with the validation issue.
func checkUser(user models.User, newUser models.User) (bool, error) {
	switch true {
	case !user.SetFirstName(newUser.First_Name):
		return true, fmt.Errorf("first name is wrongly formatted")
	case !user.SetLastName(newUser.Last_Name):
		return true, fmt.Errorf("last name is wrongly formatted")
	case !user.SetUsername(newUser.Username):
		return true, fmt.Errorf("invalid username")
	case !user.SetPassword(newUser.Password):
		return true, fmt.Errorf("invalid password")
	case !user.SetEmail(newUser.Email):
		return true, fmt.Errorf("invalid email")
	case !user.SetAddress(newUser.Address):
		return true, fmt.Errorf("invalid address")
	}
	return false, nil
}
