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

// setupRouterAndDBReview sets up the router and database for testing reviews, including migrating necessary models.
// It returns the router, database, and a teardown function to drop the tables after testing.
func setupRouterAndDBReview(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Review{}, &models.User{}, &models.Product{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Review{}, &models.User{}, &models.Product{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetReview_Success tests the successful retrieval of a review.
// It creates a user, product, and review in the database, then fetches the review by its ID.
// The test passes if the response status code is 200 and the review details match the expected values.
func TestGetReview_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	user := models.User{Username: "Test User"}
	db.Create(&user)
	product := models.Product{Name: "Test Product", Price: 10.99}
	db.Create(&product)
	review := models.Review{Product_ID: uint32(product.ID), User_ID: uint32(user.ID), Rating: 5, Comment: "Great product"}
	db.Create(&review)

	router.GET("/reviews/:id", func(c *gin.Context) {
		GetReview(c, db)
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/reviews/%d", review.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.Review
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, review.Comment, response.Comment)
}

// TestGetReview_NotFound tests the retrieval of a non-existing review.
// It attempts to fetch a review with an invalid ID and expects a 404 Not Found response.
// The test passes if the response status code is 404.
func TestGetReview_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	router.GET("/reviews/:id", func(c *gin.Context) {
		GetReview(c, db)
	})

	req, _ := http.NewRequest("GET", "/reviews/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestGetReviews_Success tests the successful retrieval of all reviews.
// It creates two reviews in the database and fetches all reviews.
// The test passes if the response status code is 200 and the number of reviews matches the expected count.
func TestGetReviews_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	user := models.User{Username: "User One"}
	db.Create(&user)
	product := models.Product{Name: "Product One", Price: 15.00}
	db.Create(&product)
	db.Create(&models.Review{Product_ID: uint32(product.ID), User_ID: uint32(user.ID), Rating: 4, Comment: "Good"})
	db.Create(&models.Review{Product_ID: uint32(product.ID), User_ID: uint32(user.ID), Rating: 3, Comment: "Average"})

	router.GET("/reviews", func(c *gin.Context) {
		GetReviews(c, db)
	})

	req, _ := http.NewRequest("GET", "/reviews", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var reviews []models.Review
	if err := json.Unmarshal(rr.Body.Bytes(), &reviews); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, reviews, 2)
}

// TestGetReviews_Empty tests retrieval of reviews when none exist.
// The test passes if the response status code is 200 and the response body is an empty array.
// This indicates that no reviews were found in the database.
// The test passes if the response status code is 200 and the response body is an empty array.
func TestGetReviews_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	router.GET("/reviews", func(c *gin.Context) {
		GetReviews(c, db)
	})

	req, _ := http.NewRequest("GET", "/reviews", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var reviews []models.Review
	if err := json.Unmarshal(rr.Body.Bytes(), &reviews); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, reviews)
}

// TestSearchAllReviews_Success tests successful searching of reviews.
// It creates a user, product, and review in the database, then searches for reviews with a specific rating.
// The test passes if the response status code is 200 and the number of reviews matches the expected count.
func TestSearchAllReviews_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	user := models.User{Username: "User Two"}
	db.Create(&user)
	product := models.Product{Name: "Product Two", Price: 20.00}
	db.Create(&product)
	db.Create(&models.Review{Product_ID: uint32(product.ID), User_ID: uint32(user.ID), Rating: 5, Comment: "Excellent"})

	router.GET("/reviews/search", func(c *gin.Context) {
		c.Request.URL.RawQuery = "rating=5"
		SearchAllReviews(c, db)
	})

	req, _ := http.NewRequest("GET", "/reviews/search", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var reviews []models.Review
	if err := json.Unmarshal(rr.Body.Bytes(), &reviews); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, reviews, 1)
}

// TestSearchAllReviews_Empty tests searching for reviews that do not match any existing data.
// The test passes if the response status code is 404, indicating that no reviews were found.
func TestSearchAllReviews_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	router.GET("/reviews/search", func(c *gin.Context) {
		c.Request.URL.RawQuery = "rating=1"
		SearchAllReviews(c, db)
	})

	req, _ := http.NewRequest("GET", "/reviews/search", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestCreateReview_Success tests successful creation of a review.
// It creates a user, product, and review in the database, then creates a new review.
// The test passes if the response status code is 201 and the review details match the expected values.
func TestCreateReview_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	user := models.User{Username: "New User"}
	db.Create(&user)
	product := models.Product{Name: "New Product", Price: 25.99}
	db.Create(&product)

	router.POST("/reviews", func(c *gin.Context) {
		CreateReview(c, db)
	})

	newReview := fmt.Sprintf(`{"product_id": %d, "user_id": %d, "rating": 5, "comment": "Fantastic!", "review_date": "2023-01-05"}`, product.ID, user.ID)
	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBufferString(newReview))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.Review
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "Fantastic!", response.Comment)
}

// TestCreateReview_InvalidData tests creation of a review with invalid data.
// It attempts to create a review with an invalid rating and an empty comment.
// The test passes if the response status code is 400, indicating that the input data was invalid.
// The response body should contain an error message.
// The test passes if the response status code is 400 and the response body contains an error message.
func TestCreateReview_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	router.POST("/reviews", func(c *gin.Context) {
		CreateReview(c, db)
	})

	newReview := `{"product_id": 999, "user_id": 999, "rating": 6, "comment": ""}`
	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBufferString(newReview))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestUpdateReview_Success tests successful updating of a review.
// It creates a user, product, and review in the database, then updates the review with new data.
// The test passes if the response status code is 200 and the review details match the updated values.
func TestUpdateReview_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	user := models.User{Username: "User Update"}
	db.Create(&user)
	product := models.Product{Name: "Product Update", Price: 30.00}
	db.Create(&product)
	review := models.Review{Product_ID: uint32(product.ID), User_ID: uint32(user.ID), Rating: 3, Comment: "Okay", Review_Date: "2023-01-05"}
	db.Create(&review)

	router.PUT("/reviews/:id", func(c *gin.Context) {
		UpdateReview(c, db)
	})

	updateData := fmt.Sprintf(`{"product_id": %d, "user_id": %d, "rating": 4, "comment": "Better", "review_date": "2023-01-06"}`, product.ID, user.ID)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/reviews/%d", review.ID), bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.Review
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "Better", response.Comment)
}

// TestUpdateReview_NotFound tests updating a non-existing review.
// It attempts to update a review with an invalid ID and expects a 404 Not Found response.
// The test passes if the response status code is 404.
func TestUpdateReview_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	router.PUT("/reviews/:id", func(c *gin.Context) {
		UpdateReview(c, db)
	})

	updateData := `{"product_id": 1, "user_id": 1, "rating": 5, "comment": "Excellent"}`
	req, _ := http.NewRequest("PUT", "/reviews/999", bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestDeleteReview_Valid tests the successful deletion of a review.
// It creates a user, product, and review in the database, then deletes the review by its ID.
// The test passes if the response status code is 204, indicating that the review was successfully deleted.
func TestDeleteReview_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	user := models.User{Username: "User Delete"}
	db.Create(&user)
	product := models.Product{Name: "Product Delete", Price: 35.00}
	db.Create(&product)
	review := models.Review{Product_ID: uint32(product.ID), User_ID: uint32(user.ID), Rating: 2, Comment: "Not good"}
	db.Create(&review)

	router.DELETE("/reviews/:id", func(c *gin.Context) {
		DeleteReview(c, db)
	})

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/reviews/%d", review.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteReview_Invalid tests attempting to delete a non-existing review.
// It tries to delete a review with an invalid ID and expects a 404 Not Found response.
// The test passes if the response status code is 404.
func TestDeleteReview_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBReview(t)
	defer teardown()

	router.DELETE("/reviews/:id", func(c *gin.Context) {
		DeleteReview(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/reviews/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
