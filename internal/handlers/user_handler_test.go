package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"E-Commerce_Website_Database/internal/models"
)

// setupRouterAndDBUser sets up the router and the database for testing user-related endpoints.
func setupRouterAndDBUser(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.User{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetUser_Success tests successful retrieval of a user by ID.
func TestGetUser_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	user := models.User{Username: "User1", Email: "user1@example.com", First_Name: "user", Last_Name: "1", Address: "123 First St"}
	db.Create(&user)

	router.GET("/users/:id", func(c *gin.Context) {
		GetUser(c, db)
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, user.Username, response.Username)
}

// TestGetUser_NotFound tests retrieval failure when a user does not exist.
func TestGetUser_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.GET("/users/:id", func(c *gin.Context) {
		GetUser(c, db)
	})

	req, _ := http.NewRequest("GET", "/users/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestGetUsers_Success tests retrieval of all users.
func TestGetUsers_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	users := []models.User{
		{Username: "user1", Email: "user1@example.com", First_Name: "First", Last_Name: "User", Address: "123 First St"},
		{Username: "user2", Email: "user2@example.com", First_Name: "Second", Last_Name: "User", Address: "456 Second St"},
	}
	for _, u := range users {
		db.Create(&u)
	}

	router.GET("/users", func(c *gin.Context) {
		GetUsers(c, db)
	})

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var responses []models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responses); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, responses, len(users))
}

// TestGetUsers_Empty tests the scenario where no users exist.
func TestGetUsers_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.GET("/users", func(c *gin.Context) {
		GetUsers(c, db)
	})

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var responses []models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responses); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, responses)
}

// TestSearchAllUsers_Success tests successful searching of users based on specific criteria.
func TestSearchAllUsers_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	user := models.User{Username: "searchuser", Email: "search@example.com", First_Name: "Search", Last_Name: "User", Address: "789 Search"}
	db.Create(&user)

	router.GET("/users/search", func(c *gin.Context) {
		SearchAllUsers(c, db)
	})

	req, _ := http.NewRequest("GET", "/users/search?username=searchuser", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var responses []models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responses); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, responses, 1)
	assert.Equal(t, "searchuser", responses[0].Username)
}

// TestSearchAllUsers_Empty tests the scenario where a search query matches no existing users.
func TestSearchAllUsers_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.GET("/users/search", func(c *gin.Context) {
		SearchAllUsers(c, db)
	})

	req, _ := http.NewRequest("GET", "/users/search?username=nonexistent", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestCreateUser_Success tests successful creation of a new user.
func TestCreateUser_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.POST("/users", func(c *gin.Context) {
		CreateUser(c, db)
	})

	newUser := `{"username": "newUser", "password": "Password123", "email": "uniq12@example.com", "first_name": "NewName", "last_name": "Username", "address": "New Street", "role": "admin"}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(newUser))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "newUser", response.Username)
}

// TestCreateUser_InvalidData tests creation of a user with invalid data.
func TestCreateUser_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.POST("/users", func(c *gin.Context) {
		CreateUser(c, db)
	})

	newUser := `{"username": "", "password": "", "email": "bademail", "first_name": "123", "last_name": "", "address": ""}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(newUser))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestUpdateUser_Success tests successful updating of an existing user.
func TestUpdateUser_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	user := models.User{Username: "updateuser", Email: "update@example.com", First_Name: "Update", Last_Name: "User", Address: "200 Update St"}
	db.Create(&user)

	router.PUT("/users/:id", func(c *gin.Context) {
		UpdateUser(c, db)
	})

	updateUser := fmt.Sprintf(`{"username": "updated", "password": "Newpassword123", "email": "updated@example.com", "first_name": "Updated", "last_name": "User", "address": "200 Updated", "role": "admin"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewBufferString(updateUser))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "updated", response.Username)
}

// TestUpdateUser_NotFound tests updating a non-existing user.
func TestUpdateUser_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.PUT("/users/:id", func(c *gin.Context) {
		UpdateUser(c, db)
	})

	updateUser := `{"username": "nonexistent", "password": "password123", "email": "nonexistent@example.com", "first_name": "Nonexistent", "last_name": "User", "address": "300 Nonexistent St"}`
	req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBufferString(updateUser))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestDeleteUser_Valid tests the successful deletion of an existing user.
func TestDeleteUser_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	user := models.User{Username: "deleteuser", Email: "delete@example.com", First_Name: "Delete", Last_Name: "User", Address: "400 Delete St"}
	db.Create(&user)

	router.DELETE("/users/:id", func(c *gin.Context) {
		DeleteUser(c, db)
	})

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", user.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteUser_Invalid tests attempting to delete a non-existing user.
func TestDeleteUser_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBUser(t)
	defer teardown()

	router.DELETE("/users/:id", func(c *gin.Context) {
		DeleteUser(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/users/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
