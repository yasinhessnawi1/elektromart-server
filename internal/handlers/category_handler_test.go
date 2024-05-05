package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"E-Commerce_Website_Database/internal/models"
)

// setupRouterAndDB sets up the router and database in memory, and returns
// a function to clean up the database after the tests.
// This function is used in the tests to set up the environment.
func setupRouterAndDBForCategoryHandler(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Category{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Category{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetCategoryIntegration checks if GetCategory function returns the correct records by given ID.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a new category in the database and sends an HTTP GET request to the /categories/:id endpoint.
// It checks the response status code, the response body, and the returned category data.
func TestGetCategoryIntegration(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Test Category", Description: "Test Description"})

	router.GET("/categories/:id", func(c *gin.Context) {
		GetCategory(c, db)
	})

	req, _ := http.NewRequest("GET", "/categories/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, float64(1), response["ID"])
	assert.Equal(t, "Test Category", response["name"])
	assert.Equal(t, "Test Description", response["description"])

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestGetCategoryIntegrationInvalid checks if GetCategory function returns the correct error.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a new category in the database and sends an HTTP GET request to the /categories/:id endpoint with an invalid ID.
// It checks the response status code and the error message in the response body.
func TestGetCategoryIntegrationInvalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Test Category", Description: "Test Description"})

	router.GET("/categories/:id", func(c *gin.Context) {
		GetCategory(c, db)
	})

	req, _ := http.NewRequest("GET", "/categories/2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, response["error"], "Category not found")
}

// TestGetCategoriesIntegration checks if GetCategories function returns all categories from the database.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates two new categories in the database and sends an HTTP GET request to the /categories endpoint.
// It checks the response status code, the response body, and the returned categories data.
func TestGetCategoriesIntegration(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Category 1", Description: "Description 1"})
	db.Create(&models.Category{Model: gorm.Model{ID: 2}, Name: "Category 2", Description: "Description 2"})

	router.GET("/categories", func(c *gin.Context) {
		GetCategories(c, db)
	})

	req, _ := http.NewRequest("GET", "/categories", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Category
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if all categories are retrieved
	assert.Equal(t, 2, len(response))
	assert.Equal(t, "Category 1", response[0].Name)
	assert.Equal(t, "Description 1", response[0].Description)
	assert.Equal(t, "Category 2", response[1].Name)
	assert.Equal(t, "Description 2", response[1].Description)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSearchAllCategoriesIntegration checks if SearchAllCategories function returns categories based on search parameters.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates two new categories in the database and sends an HTTP GET request to the /categories/search endpoint with a search parameter.
// It checks the response status code, the response body, and the returned categories data.
func TestSearchAllCategoriesIntegration(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Category 1", Description: "Description 1"})
	db.Create(&models.Category{Model: gorm.Model{ID: 2}, Name: "Category 2", Description: "Description 2"})

	router.GET("/categories/search", func(c *gin.Context) {
		SearchAllCategories(c, db)
	})

	req, _ := http.NewRequest("GET", "/categories/search?name=Category 1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Category
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if the correct category is retrieved
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Category 1", response[0].Name)
	assert.Equal(t, "Description 1", response[0].Description)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSearchAllCategoriesIntegrationEmpty checks if SearchAllCategories function returns a message when no categories match the search criteria.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a new category in the database and sends an HTTP GET request to the /categories/search endpoint with a search parameter that does not match any categories.
// It checks the response status code and the error message in the response body.
func TestSearchAllCategoriesIntegrationEmpty(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Category 1", Description: "Description 1"})

	router.GET("/categories/search", func(c *gin.Context) {
		SearchAllCategories(c, db)
	})

	req, _ := http.NewRequest("GET", "/categories/search?name=Category 2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if the response contains the appropriate message
	assert.Equal(t, "No category found", response["error"])

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestCreateCategory_Success ensures that a category can be successfully created with valid data.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It sends an HTTP POST request to the /categories endpoint with valid category data.
// It checks the response status code, the response body, and the returned category data.
func TestCreateCategory_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	router.POST("/categories", func(c *gin.Context) {
		CreateCategory(c, db)
	})

	newCategory := `{"name":"Valid Category", "description":"Valid Description"}`
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(newCategory))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check response body is correct
	assert.Equal(t, "Valid Category", response["name"])
	assert.Equal(t, "Valid Description", response["description"])
}

// TestCreateCategory_InvalidData checks the response when incomplete or incorrect data is sent.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It sends an HTTP POST request to the /categories endpoint with invalid category data.
// It checks the response status code and the error message in the response body.
func TestCreateCategory_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	router.POST("/categories", func(c *gin.Context) {
		CreateCategory(c, db)
	})

	newCategory := `{"name":""}`
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(newCategory))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check response body
	assert.Contains(t, response["error"], "Validation error")
}

// TestUpdateCategoryValid checks the ability to update an existing category.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a category in the database and sends an HTTP PUT request to the /categories/:id endpoint with updated data.
// It checks the response status code, the response body, and the updated category data.
func TestUpdateCategoryValid(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Old Category", Description: "Old Description"})

	// Set up the PUT route
	router.PUT("/categories/:id", func(c *gin.Context) {
		UpdateCategory(c, db)
	})

	// Update category via HTTP PUT
	updatedCategory := `{"name":"Updated Category", "description":"Updated Description"}`
	req, _ := http.NewRequest("PUT", "/categories/1", strings.NewReader(updatedCategory))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check that the response contains the updated data
	assert.Equal(t, "Updated Category", response["name"])
	assert.Equal(t, "Updated Description", response["description"])
}

// TestUpdateCategoryInvalid checks the error message with invalid data.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a category in the database and sends an HTTP PUT request to the /categories/:id endpoint with invalid data.
// It checks the response status code and the error message in the response body.
func TestUpdateCategoryInvalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Old Category", Description: "Old Description"})

	// Set up the PUT route
	router.PUT("/categories/:id", func(c *gin.Context) {
		UpdateCategory(c, db)
	})

	// Update category via HTTP PUT
	updatedCategory := `{"name":"", "description":""}`
	req, _ := http.NewRequest("PUT", "/categories/1", strings.NewReader(updatedCategory))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check that the response contains the correct error message
	assert.Contains(t, response["error"], "Validation error")
}

// TestDeleteCategoryValid checks that a category is deleted from the database.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a category in the database and sends an HTTP DELETE request to the /categories/:id endpoint.
// If the category is successfully deleted, the response status code should be 204 No Content.
// If the category is not found, the response status code should be 404 Not Found.
func TestDeleteCategoryValid(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Category to Delete", Description: "Description"})

	// Set up the DELETE route
	router.DELETE("/categories/:id", func(c *gin.Context) {
		DeleteCategory(c, db)
	})

	// Delete category via HTTP DELETE
	req, _ := http.NewRequest("DELETE", "/categories/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is correct
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteCategoryInvalid checks the delete category with invalid ID.
// It uses the setupRouterAndDBForCategoryHandler function to set up the environment.
// It creates a category in the database and sends an HTTP DELETE request to the /categories/:id endpoint with an invalid ID.
// It checks the response status code and the error message in the response body.
func TestDeleteCategoryInvalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBForCategoryHandler(t)
	defer teardown()

	db.Create(&models.Category{Model: gorm.Model{ID: 1}, Name: "Category to Delete", Description: "Description"})

	// Set up the DELETE route
	router.DELETE("/categories/:id", func(c *gin.Context) {
		DeleteCategory(c, db)
	})

	// Delete category via HTTP DELETE
	req, _ := http.NewRequest("DELETE", "/categories/2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if the response contains the error message
	assert.Contains(t, response["error"], "Category not found")
}
