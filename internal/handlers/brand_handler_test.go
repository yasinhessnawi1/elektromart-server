package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"E-Commerce_Website_Database/internal/models"
)

// setupRouterAndDB sets up the router and database in memory, and returns a function to clean up the database after the tests.
// It returns the router, database, and a teardown function.
func setupRouterAndDB(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Brands{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Brands{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetBrandIntegration checks at this function and return the correct records by the given ID.
// It should return a status code of 200 and a JSON response with the brand details.
func TestGetBrandIntegration(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Test Brand", Description: "Test Description"})

	router.GET("/brands/:id", func(c *gin.Context) {
		GetBrand(c, db)
	})

	req, _ := http.NewRequest("GET", "/brands/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, float64(1), response["ID"])
	assert.Equal(t, "Test Brand", response["name"])
	assert.Equal(t, "Test Description", response["description"])

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestGetBrandIntegration checks at this function and return the correct records by the given ID.
// It checks the response when the ID is invalid.
// It should return an error message and a status code of 404.
// The error message should indicate that the brand was not found.
func TestGetBrandIntegrationInvalid(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Test Brand", Description: "Test Description"})

	router.GET("/brands/:id", func(c *gin.Context) {
		GetBrand(c, db)
	})

	req, _ := http.NewRequest("GET", "/brands/2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, response["error"], "Brand not found")
}

// TestGetBrandsIntegration checks if the GetBrands function returns all brands from the database.
// It should return a status code of 200 and a JSON response with all brands.
// The response should contain the correct brand details.
func TestGetBrandsIntegration(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Brand 1", Description: "Description 1"})
	db.Create(&models.Brands{Model: gorm.Model{ID: 2}, Name: "Brand 2", Description: "Description 2"})

	router.GET("/brands", func(c *gin.Context) {
		GetBrands(c, db)
	})

	req, _ := http.NewRequest("GET", "/brands", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Brands
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if all brands are retrieved
	assert.Equal(t, 2, len(response))
	assert.Equal(t, "Brand 1", response[0].Name)
	assert.Equal(t, "Description 1", response[0].Description)
	assert.Equal(t, "Brand 2", response[1].Name)
	assert.Equal(t, "Description 2", response[1].Description)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSearchAllBrandsIntegration checks if the SearchAllBrands function returns brands based on search parameters.
// It should return a status code of 200 and a JSON response with the matching brands.
// The response should contain the correct brand details.
func TestSearchAllBrandsIntegration(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Brand 1", Description: "Description 1"})
	db.Create(&models.Brands{Model: gorm.Model{ID: 2}, Name: "Brand 2", Description: "Description 2"})

	router.GET("/brands/search", func(c *gin.Context) {
		SearchAllBrands(c, db)
	})

	req, _ := http.NewRequest("GET", "/brands/search?name=Brand 1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.Brands
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if the correct brand is retrieved
	assert.NotNil(t, response)
	assert.Equal(t, "Brand 1", response.Name)
	assert.Equal(t, "Description 1", response.Description)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSearchAllBrandsIntegrationEmpty checks if the SearchAllBrands function returns
// a message when no brands match the search criteria.
// It should return a status code of 404 and a JSON response with an error message.
func TestSearchAllBrandsIntegrationEmpty(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Brand 1", Description: "Description 1"})

	router.GET("/brands/search", func(c *gin.Context) {
		SearchAllBrands(c, db)
	})

	req, _ := http.NewRequest("GET", "/brands/search?name=Brand 2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.Brands
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check the status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// TestCreateBrand_Success ensures that a brand can be successfully created with valid data.
// It checks the response body and status code.
// It should return a status code of 201 and a JSON response with the created brand details.
func TestCreateBrand_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	router.POST("/brands", func(c *gin.Context) {
		CreateBrand(c, db)
	})

	newBrand := `{"name":"Valid Brand", "description":"Valid Description"}`
	req, _ := http.NewRequest("POST", "/brands", bytes.NewBufferString(newBrand))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check response body is correct
	assert.Equal(t, "Valid Brand", response["name"])
	assert.Equal(t, "Valid Description", response["description"])
}

// TestCreateBrand_InvalidData checks the response when incomplete or incorrect data is sent.
// It should return an error message and a status code of 400.
func TestCreateBrand_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	router.POST("/brands", func(c *gin.Context) {
		CreateBrand(c, db)
	})

	newBrand := `{"name":""}`
	req, _ := http.NewRequest("POST", "/brands", bytes.NewBufferString(newBrand))
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

// TestUpdateBrandValid Checks the ability to update an existing brand with valid data.
// It should return a status code of 200 and a JSON response with the updated brand details.
// The response should contain the correct brand details.
// The brand should be updated in the database.
// The response should contain the updated brand details.
func TestUpdateBrandValid(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Old Brand", Description: "Old Description"})

	// Set up the PUT route
	router.PUT("/brands/:id", func(c *gin.Context) {
		UpdateBrand(c, db)
	})

	// Update brand via HTTP PUT
	updatedBrand := `{"name":"Updated Brand", "description":"Updated Description"}`
	req, _ := http.NewRequest("PUT", "/brands/1", strings.NewReader(updatedBrand))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check that the response contains the updated data
	assert.Equal(t, "Updated Brand", response["name"])
	assert.Equal(t, "Updated Description", response["description"])
}

// TestUpdateBrandInvalid Checks the error message with invalid data.
// It should return an error message and a status code of 400.
func TestUpdateBrandInvalid(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Old Brand", Description: "Old Description"})

	// Set up the PUT route
	router.PUT("/brands/:id", func(c *gin.Context) {
		UpdateBrand(c, db)
	})

	// Update brand via HTTP PUT
	updatedBrand := `{"name":"", "description":""}`
	req, _ := http.NewRequest("PUT", "/brands/1", strings.NewReader(updatedBrand))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// checks that the response contains the correct error message
	assert.Contains(t, response["error"], "Validation error")

}

// TestDeleteBrandValid checks that a brand is deleted from the database
func TestDeleteBrandValid(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Brand to Delete", Description: "Description"})

	// Set up the DELETE route
	router.DELETE("/brands/:id", func(c *gin.Context) {
		DeleteBrand(c, db)
	})

	// Delete brand via HTTP DELETE
	req, _ := http.NewRequest("DELETE", "/brands/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is correct
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteBrandInvalid checks the delete brand with invalid ID and return an error message.
// It should return an error message and a status code of 404.
// The error message should indicate that the brand was not found.
func TestDeleteBrandInvalid(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	db.Create(&models.Brands{Model: gorm.Model{ID: 1}, Name: "Brand to Delete", Description: "Description"})

	// Set up the DELETE route
	router.DELETE("/brands/:id", func(c *gin.Context) {
		DeleteBrand(c, db)
	})

	// Delete brand via HTTP DELETE
	req, _ := http.NewRequest("DELETE", "/brands/2", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	// Check if the response contains the error message
	assert.Contains(t, response["error"], "Brands not found")
}
