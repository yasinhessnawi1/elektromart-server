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

// setupRouterAndDBShippingDetail sets up the router and the database for testing, including the Order model.
// It returns the router, database, and a function to clean up the database after testing.
func setupRouterAndDBShippingDetail(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.ShippingDetails{}, &models.Order{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.ShippingDetails{}, &models.Order{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetShippingDetail_Success tests successful retrieval of a shipping detail by ID.
// It creates a new shipping detail, then attempts to retrieve it by ID.
// If the retrieval is successful, it checks the response status code and the shipping detail's address.
// If the retrieval fails, it fails the test.
func TestGetShippingDetail_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	order := models.Order{Total_amount: 150.50}
	db.Create(&order)

	shippingDetail := models.ShippingDetails{Order_ID: uint32(order.ID), Address: "123 First St", Shipping_Date: "2023-04-01", Estimated_Arrival: "2023-04-05", Status: "shipped"}
	db.Create(&shippingDetail)

	router.GET("/shippingDetails/:id", func(c *gin.Context) {
		GetShippingDetail(c, db)
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/shippingDetails/%d", shippingDetail.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.ShippingDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, shippingDetail.Address, response.Address)
}

// TestGetShippingDetail_NotFound tests retrieval failure when a shipping detail does not exist.
// It attempts to retrieve a non-existing shipping detail by ID.
// If the retrieval fails, it checks the response status code.
// If the retrieval is successful, it fails the test.
func TestGetShippingDetail_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	router.GET("/shippingDetails/:id", func(c *gin.Context) {
		GetShippingDetail(c, db)
	})

	req, _ := http.NewRequest("GET", "/shippingDetails/9999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestGetShippingDetails_Success tests retrieval of all shipping details.
// It creates multiple shipping details, then attempts to retrieve all of them.
// If the retrieval is successful, it checks the response status code and the number of shipping details.
// If the retrieval fails, it fails the test.
func TestGetShippingDetails_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	order := models.Order{Total_amount: 200.00}
	db.Create(&order)

	shippingDetails := []models.ShippingDetails{
		{Order_ID: uint32(order.ID), Address: "123 First St", Shipping_Date: "2023-04-01", Estimated_Arrival: "2023-04-05", Status: "shipped"},
		{Order_ID: uint32(order.ID), Address: "456 Second St", Shipping_Date: "2023-05-01", Estimated_Arrival: "2023-05-05", Status: "pending"},
	}
	for _, detail := range shippingDetails {
		db.Create(&detail)
	}

	router.GET("/shippingDetails", func(c *gin.Context) {
		GetShippingDetails(c, db)
	})

	req, _ := http.NewRequest("GET", "/shippingDetails", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var responses []models.ShippingDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &responses); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, responses, len(shippingDetails))
}

// TestGetShippingDetails_Empty tests the scenario where no shipping details exist.
// It attempts to retrieve all shipping details when none exist.
// If the retrieval is successful, it checks the response status code and the number of shipping details.
// If the retrieval fails, it fails the test.
func TestGetShippingDetails_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	router.GET("/shippingDetails", func(c *gin.Context) {
		GetShippingDetails(c, db)
	})

	req, _ := http.NewRequest("GET", "/shippingDetails", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var responses []models.ShippingDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &responses); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, responses)
}

// TestSearchAllShippingDetails_Success tests the successful search for shipping details based on specific criteria.
// It creates a shipping detail and attempts to search for it based on the status field.
// If the search is successful, it checks the response status code and the number of matching shipping details.
// If the search fails, it fails the test.
func TestSearchAllShippingDetails_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	order := models.Order{Total_amount: 250.50}
	db.Create(&order)

	shippingDetail := models.ShippingDetails{Order_ID: uint32(order.ID), Address: "789 Off St", Shipping_Date: "2023-06-01", Estimated_Arrival: "2023-06-05", Status: "in transit"}
	db.Create(&shippingDetail)

	router.GET("/shippingDetails/search/", func(c *gin.Context) {
		SearchAllShippingDetails(c, db)
	})

	req, _ := http.NewRequest("GET", "/shippingDetails/search/?status=in+transit", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var responses []models.ShippingDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &responses); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, responses, 1)
	assert.Equal(t, "789 Off St", responses[0].Address)
}

// TestSearchAllShippingDetails_Empty tests the scenario where a search query matches no existing shipping details.
// It attempts to search for shipping details based on a non-existing status.
// If the search is successful, it checks the response status code.
// If the search fails, it fails the test.
func TestSearchAllShippingDetails_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	router.GET("/shippingDetails/search/", func(c *gin.Context) {
		SearchAllShippingDetails(c, db)
	})

	req, _ := http.NewRequest("GET", "/shippingDetails/search/?status=delivered", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestCreateShippingDetail_Success tests the successful creation of a new shipping detail.
// It creates a new shipping detail and attempts to create it in the database.
// If the creation is successful, it checks the response status code and the shipping detail's address.
// If the creation fails, it fails the test.
func TestCreateShippingDetail_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	order := models.Order{Total_amount: 300.50}
	db.Create(&order)

	router.POST("/shippingDetails", func(c *gin.Context) {
		CreateShippingDetail(c, db)
	})

	newDetail := `{"order_id": %d, "address": "New Address", "shipping_date": "2023-07-01", "estimated_arrival": "2023-07-05", "status": "pending"}`
	newDetail = fmt.Sprintf(newDetail, order.ID)
	req, _ := http.NewRequest("POST", "/shippingDetails", bytes.NewBufferString(newDetail))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.ShippingDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "New Address", response.Address)
}

// TestCreateShippingDetail_InvalidData tests the creation of a shipping detail with invalid data.
// It attempts to create a shipping detail with invalid data.
// If the creation fails, it checks the response status code.
// If the creation is successful, it fails the test.
func TestCreateShippingDetail_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	router.POST("/shippingDetails", func(c *gin.Context) {
		CreateShippingDetail(c, db)
	})

	newDetail := `{"order_id": 999, "address": "", "shipping_date": "bad date", "estimated_arrival": "", "status": "unknown"}`
	req, _ := http.NewRequest("POST", "/shippingDetails", bytes.NewBufferString(newDetail))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// TestUpdateShippingDetail_Success tests the successful update of an existing shipping detail.
// It creates an original shipping detail, then attempts to update it with new data.
// If the update is successful, it checks the response status code and the updated shipping detail's address.
// If the update fails, it fails the test.
func TestUpdateShippingDetail_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	order := models.Order{Total_amount: 350.50}
	db.Create(&order)
	originalDetail := models.ShippingDetails{Order_ID: uint32(order.ID), Address: "Original Address", Shipping_Date: "2023-08-01", Estimated_Arrival: "2023-08-05", Status: "pending"}
	db.Create(&originalDetail)

	router.PUT("/shippingDetails/:id", func(c *gin.Context) {
		UpdateShippingDetail(c, db)
	})

	updateDetail := fmt.Sprintf(`{"order_id": %d, "address": "Updated Address", "shipping_date": "2023-08-01", "estimated_arrival": "2023-08-10", "status": "shipped"}`, order.ID)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/shippingDetails/%d", originalDetail.ID), bytes.NewBufferString(updateDetail))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.ShippingDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "Updated Address", response.Address)
	assert.Equal(t, "2023-08-10", response.Estimated_Arrival)
}

// TestUpdateShippingDetail_NotFound tests the scenario where an attempt is made to update a non-existing shipping detail.
// It attempts to update a non-existing shipping detail with new data.
// If the update fails, it checks the response status code.
// If the update is successful, it fails the test.
func TestUpdateShippingDetail_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	router.PUT("/shippingDetails/:id", func(c *gin.Context) {
		UpdateShippingDetail(c, db)
	})

	updateDetail := `{"order_id": 1, "address": "Nonexistent Address", "shipping_date": "2023-09-01", "estimated_arrival": "2023-09-05", "status": "pending"}`
	req, _ := http.NewRequest("PUT", "/shippingDetails/999", bytes.NewBufferString(updateDetail))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestDeleteShippingDetail_Valid tests the successful deletion of an existing shipping detail.
// It creates a shipping detail, then attempts to delete it.
// If the deletion is successful, it checks the response status code.
// If the deletion fails, it fails the test.
func TestDeleteShippingDetail_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	order := models.Order{Total_amount: 400.50}
	db.Create(&order)
	detailToDelete := models.ShippingDetails{Order_ID: uint32(order.ID), Address: "Delete Me", Shipping_Date: "2023-10-01", Estimated_Arrival: "2023-10-05", Status: "pending"}
	db.Create(&detailToDelete)

	router.DELETE("/shippingDetails/:id", func(c *gin.Context) {
		DeleteShippingDetail(c, db)
	})

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/shippingDetails/%d", detailToDelete.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteShippingDetail_Invalid tests the scenario where an attempt is made to delete a non-existing shipping detail.
// It attempts to delete a non-existing shipping detail.
// If the deletion fails, it checks the response status code.
// If the deletion is successful, it fails the test.
func TestDeleteShippingDetail_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBShippingDetail(t)
	defer teardown()

	router.DELETE("/shippingDetails/:id", func(c *gin.Context) {
		DeleteShippingDetail(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/shippingDetails/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
