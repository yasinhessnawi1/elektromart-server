package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// GetUser retrieves a single user by ID from the URL parameters.
// It returns the user details or an error message if the user is not found.
// If the user is found, it responds with an HTTP 200 OK status and the user details in JSON format.
// If the user is not found, it responds with an HTTP 404 Not Found status.
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
// If there are no users in the database, it responds with an HTTP 404 Not Found status.
// If the retrieval is successful, it responds with an HTTP 200 OK status and the list of users in JSON format.
// If there is an error during retrieval, it responds with an HTTP 500 Internal Server Error status.
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
// If no users are found, it responds with an HTTP 404 Not Found status.
// If the search is successful, it responds with an HTTP 200 OK status and the list of users in JSON format.
// If there is an error during retrieval, it responds with an HTTP 500 Internal Server Error status.
func SearchAllUsers(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"username", "email", "first_name", "last_name", "address"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			searchParams[field] = cleanValue
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
// It validates the input and stores the new user in the database, responding to the created user or an error message.
// If the user is created successfully, it responds with an HTTP 201 Created status and the user details in JSON format.
// If there is an error during creation, it responds with an HTTP 500 Internal Server Error status.
func CreateUser(c *gin.Context, db *gorm.DB) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}

	user := models.User{
		Username:   newUser.Username,
		Password:   string(hashedPassword),
		Email:      newUser.Email,
		First_Name: newUser.First_Name,
		Last_Name:  newUser.Last_Name,
		Address:    newUser.Address,
		Mobile:     newUser.Mobile,
		Role:       "regular",
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkUser(user, newUser, true); failed {
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
// If the user is not found, it responds with an HTTP 404 Not Found status.
// If the input data is invalid, it responds with an HTTP 400 Bad Request status.
// If the update is successful, it responds with an HTTP 200 OK status and the updated user details in JSON format.
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

	// Load the existing user
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found during update"})
		return
	}
	if newUser.Password != "" {
		if tools.CheckPassword(newUser.Password) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
				return
			}
			user.Password = string(hashedPassword)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not valid, must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number and one special character"})
			return

		}
	}

	// Update user fields
	user.Username = newUser.Username
	user.Email = newUser.Email
	user.First_Name = newUser.First_Name
	user.Last_Name = newUser.Last_Name
	user.Address = newUser.Address
	user.Mobile = newUser.Mobile

	if failed, err := checkUser(user, newUser, false); failed && err != nil {
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
// If the user is not found, it responds with an HTTP 404 Not Found status.
// If the deletion is successful, it responds with an HTTP 204 No Content status.
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
// If the data is valid, it returns false and nil.
// If the data is invalid, it returns true and an error message.
func checkUser(user models.User, newUser models.User, isCreating bool) (bool, error) {
	switch true {
	case !user.SetFirstName(newUser.First_Name):
		return true, fmt.Errorf("first name is wrong formatted")
	case !user.SetLastName(newUser.Last_Name):
		return true, fmt.Errorf("last name is wrong formatted")
	case !user.SetUsername(newUser.Username):
		return true, fmt.Errorf("invalid username")
	case !user.SetPassword(newUser.Password):
		if isCreating {
			return true, fmt.Errorf("invalid password")
		}
		return false, nil
	case !user.SetEmail(newUser.Email):
		return true, fmt.Errorf("invalid email")
	case !user.SetAddress(newUser.Address):
		return true, fmt.Errorf("invalid address")
	case !user.SetPhone(newUser.Mobile):
		return true, fmt.Errorf("invalid mobile")
	}
	return false, nil
}

// GetUserByUN retrieves a single user by username.
// It returns the user details or an error message if the user is not found.
// If the user is found, it responds with an HTTP 200 OK status and the user details in JSON format.
// If the user is not found, it responds with an HTTP 404 Not Found status.
func GetUserByUN(username string, db *gorm.DB) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
