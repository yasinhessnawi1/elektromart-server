package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// setupRouterAndDBOrder sets up the router and database in memory, and returns a function to clean up the database after tests.
func setupRouterAndDBOrder(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Migrate both the Order and User models
	if err := db.AutoMigrate(&models.Order{}, &models.User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Order{}, &models.User{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetOrder_Success checks if GetOrder returns the correct order by given ID.
func TestGetOrder_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Insert mock order
	order := models.Order{User_ID: 1, Order_date: "2021-09-15", Total_amount: 100.00, Status: "completed"}
	db.Create(&order)

	router.GET("/orders/:id", func(c *gin.Context) {
		GetOrder(c, db)
	})

	// Create a request to get the order
	orderID := strconv.Itoa(int(order.ID))
	req, _ := http.NewRequest("GET", "/orders/"+orderID, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, order.ID, response.ID)
	assert.Equal(t, order.User_ID, response.User_ID)
	assert.Equal(t, order.Order_date, response.Order_date)
	assert.Equal(t, order.Total_amount, response.Total_amount)
	assert.Equal(t, order.Status, response.Status)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestGetOrder_NotFound checks the response for a non-existent order.
func TestGetOrder_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	router.GET("/orders/:id", func(c *gin.Context) {
		GetOrder(c, db)
	})

	req, _ := http.NewRequest("GET", "/orders/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Order not found")
}

// TestGetOrders_Success checks if GetOrders returns all orders from the database.
func TestGetOrders_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Insert mock orders
	db.Create(&models.Order{User_ID: 1, Order_date: "2021-09-15", Total_amount: 100.00, Status: "completed"})
	db.Create(&models.Order{User_ID: 2, Order_date: "2021-09-16", Total_amount: 200.00, Status: "pending"})

	router.GET("/orders", func(c *gin.Context) {
		GetOrders(c, db)
	})

	req, _ := http.NewRequest("GET", "/orders", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, 2, len(response))
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestGetOrders_Empty checks the response when no orders are available in the database.
func TestGetOrders_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	router.GET("/orders", func(c *gin.Context) {
		GetOrders(c, db)
	})

	req, _ := http.NewRequest("GET", "/orders", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response []models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}
	assert.Equal(t, 0, len(response))
}

// TestSearchAllOrders_Success checks if SearchAllOrders returns orders based on search parameters.
func TestSearchAllOrders_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Insert mock orders
	db.Create(&models.Order{User_ID: 1, Order_date: "2021-09-15", Total_amount: 100.00, Status: "completed"})
	db.Create(&models.Order{User_ID: 2, Order_date: "2021-09-16", Total_amount: 200.00, Status: "pending"})

	router.GET("/orders/search", func(c *gin.Context) {
		SearchAllOrders(c, db)
	})

	req, _ := http.NewRequest("GET", "/orders/search?status=completed", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, 1, len(response))
	assert.Equal(t, "completed", response[0].Status)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSearchAllOrders_NotFound checks if SearchAllOrders responds correctly when no orders match the search criteria.
func TestSearchAllOrders_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	router.GET("/orders/search", func(c *gin.Context) {
		SearchAllOrders(c, db)
	})

	req, _ := http.NewRequest("GET", "/orders/search?status=shipped", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "No orders found", response["error"])
}

// TestCreateOrder_Success checks that an order can be successfully created with valid data.
func TestCreateOrder_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Create a user first
	user := models.User{Username: "testuser", Password: "testpassword"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	router.POST("/orders", func(c *gin.Context) {
		CreateOrder(c, db)
	})

	newOrder := `{"user_id": ` + strconv.Itoa(int(user.ID)) + `, "order_date": "2021-09-15", "total_amount": 100.00, "status": "completed"}`
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString(newOrder))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, uint32(user.ID), response.User_ID)
	assert.Equal(t, "2021-09-15", response.Order_date)
	assert.Equal(t, 100.00, response.Total_amount)
	assert.Equal(t, "completed", response.Status)
}

// TestCreateOrder_InvalidData checks the response when incomplete or incorrect data is sent.
func TestCreateOrder_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	router.POST("/orders", func(c *gin.Context) {
		CreateOrder(c, db)
	})

	newOrder := `{"user_id": "", "order_date": "", "total_amount": "100.00", "status": "completed"}`
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString(newOrder))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Invalid JSON data")
}

// TestUpdateOrder_Valid checks the ability to update an existing order.
func TestUpdateOrder_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Create a user first
	user := models.User{Username: "testuser", Password: "testpassword"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Create original order with the user's actual ID
	order := models.Order{User_ID: uint32(user.ID), Order_date: "2021-09-15", Total_amount: 100.00, Status: "completed"}
	db.Create(&order)

	router.PUT("/orders/:id", func(c *gin.Context) {
		UpdateOrder(c, db)
	})

	// Ensure you're using the correct User_ID
	updatedOrder := `{"user_id": ` + strconv.Itoa(int(user.ID)) + `, "order_date": "2021-10-15", "total_amount": 150.00, "status": "pending"}`
	orderID := strconv.Itoa(int(order.ID))
	req, _ := http.NewRequest("PUT", "/orders/"+orderID, bytes.NewBufferString(updatedOrder))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.Order
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, uint32(user.ID), response.User_ID)
	assert.Equal(t, "2021-10-15", response.Order_date)
	assert.Equal(t, 150.00, response.Total_amount)
	assert.Equal(t, "pending", response.Status)
}

// TestUpdateOrder_Invalid checks the error message with invalid data.
func TestUpdateOrder_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Create original order
	order := models.Order{User_ID: 1, Order_date: "2021-09-15", Total_amount: 100.00, Status: "completed"}
	db.Create(&order)

	router.PUT("/orders/:id", func(c *gin.Context) {
		UpdateOrder(c, db)
	})

	// Update order via HTTP PUT with invalid data
	updatedOrder := `{"user_id": "", "order_date": "", "total_amount": "", "status": "pending"}`
	orderID := strconv.Itoa(int(order.ID))
	req, _ := http.NewRequest("PUT", "/orders/"+orderID, bytes.NewBufferString(updatedOrder))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Invalid JSON data")
}

// TestDeleteOrder_Valid checks that an order is deleted from the database.
func TestDeleteOrder_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrder(t)
	defer teardown()

	// Create an order to delete
	order := models.Order{User_ID: 1, Order_date: "2021-09-15", Total_amount: 100.00, Status: "completed"}
	db.Create(&order)

	router.DELETE("/orders/:id", func(c *gin.Context) {
		DeleteOrder(c, db)
	})

	// Delete order via HTTP DELETE
	orderID := strconv.Itoa(int(order.ID))
	req, _ := http.NewRequest("DELETE", "/orders/"+orderID, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteOrder_Invalid checks the delete order with invalid ID.
func TestDeleteOrder_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDB(t)
	defer teardown()

	router.DELETE("/orders/:id", func(c *gin.Context) {
		DeleteOrder(c, db)
	})

	// Delete order via HTTP DELETE with non-existing ID
	req, _ := http.NewRequest("DELETE", "/orders/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Order not found")
}
