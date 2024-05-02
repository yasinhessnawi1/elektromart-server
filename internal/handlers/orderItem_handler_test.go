package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// setupRouterAndDBOrderItem sets up the router and database in memory, including the migration of Order, Product, and OrderItem models.
func setupRouterAndDBOrderItem(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Order{}, &models.Product{}, &models.OrderItem{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Order{}, &models.Product{}, &models.OrderItem{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetOrderItem_Success verifies that fetching an existing order item by ID correctly returns the order item.
func TestGetOrderItem_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	// Create supporting records and an order item
	order := models.Order{Total_amount: 100.00}
	product := models.Product{Name: "Test Product", Price: 10.00}
	db.Create(&order)
	db.Create(&product)
	orderItem := models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 5, Subtotal: 50.00}
	db.Create(&orderItem)

	router.GET("/orderItems/:id", func(c *gin.Context) {
		GetOrderItem(c, db)
	})

	req, _ := http.NewRequest("GET", "/orderItems/"+strconv.Itoa(int(orderItem.ID)), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.OrderItem
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, orderItem.ID, response.ID)
	assert.Equal(t, orderItem.Quantity, response.Quantity)
}

// TestGetOrderItem_NotFound checks the response for a non-existent order item.
func TestGetOrderItem_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	router.GET("/orderItems/:id", func(c *gin.Context) {
		GetOrderItem(c, db)
	})

	req, _ := http.NewRequest("GET", "/orderItems/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Order item not found")
}

// TestGetOrderItems_Success checks if GetOrderItems returns all order items from the database.
func TestGetOrderItems_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	order := models.Order{Total_amount: 200.00}
	product := models.Product{Name: "Test Product", Price: 20.00}
	db.Create(&order)
	db.Create(&product)
	db.Create(&models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 2, Subtotal: 40.00})
	db.Create(&models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 3, Subtotal: 60.00})

	router.GET("/orderItems", func(c *gin.Context) {
		GetOrderItems(c, db)
	})

	req, _ := http.NewRequest("GET", "/orderItems", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.OrderItem
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 2, len(response))
}

// TestGetOrderItems_Empty checks the response when no order items are available in the database.
func TestGetOrderItems_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	router.GET("/orderItems", func(c *gin.Context) {
		GetOrderItems(c, db)
	})

	req, _ := http.NewRequest("GET", "/orderItems", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response []models.OrderItem
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, 0, len(response))
}

// TestSearchAllOrderItems_Success checks if SearchAllOrderItems returns order items based on search parameters.
func TestSearchAllOrderItems_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	order := models.Order{Total_amount: 200.00}
	product := models.Product{Name: "Test Product", Price: 20.00}
	db.Create(&order)
	db.Create(&product)
	db.Create(&models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 10, Subtotal: 200.00})

	router.GET("/orderItems/search", func(c *gin.Context) {
		SearchAllOrderItems(c, db)
	})

	req, _ := http.NewRequest("GET", "/orderItems/search?quantity=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.OrderItem
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, response, 1)
	assert.Equal(t, 10, response[0].Quantity)
}

// TestSearchAllOrderItems_NotFound checks if SearchAllOrderItems responds correctly when no order items match the search criteria.
func TestSearchAllOrderItems_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	router.GET("/orderItems/search", func(c *gin.Context) {
		SearchAllOrderItems(c, db)
	})

	req, _ := http.NewRequest("GET", "/orderItems/search?quantity=100", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "No order items found")
}

// TestCreateOrderItem_Success checks that an order item can be successfully created with valid data.
func TestCreateOrderItem_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	order := models.Order{Total_amount: 100.00}
	product := models.Product{Name: "Test Product", Price: 10.00}
	db.Create(&order)
	db.Create(&product)

	router.POST("/orderItems", func(c *gin.Context) {
		CreateOrderItem(c, db)
	})

	newOrderItem := fmt.Sprintf(`{"order_id": %d, "product_id": %d, "quantity": 5, "subtotal": 50.00}`, order.ID, product.ID)
	req, _ := http.NewRequest("POST", "/orderItems", bytes.NewBufferString(newOrderItem))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.OrderItem
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, uint32(order.ID), response.Order_ID)
	assert.Equal(t, uint32(product.ID), response.Product_ID)
	assert.Equal(t, 5, response.Quantity)
	assert.Equal(t, 50.00, response.Subtotal)
}

// TestCreateOrderItem_InvalidData checks the response when incomplete or incorrect data is sent.
func TestCreateOrderItem_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	router.POST("/orderItems", func(c *gin.Context) {
		CreateOrderItem(c, db)
	})

	newOrderItem := `{"order_id": "", "product_id": "", "quantity": "five", "subtotal": "fifty"}`
	req, _ := http.NewRequest("POST", "/orderItems", bytes.NewBufferString(newOrderItem))
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

// TestUpdateOrderItem_Valid checks the ability to update an existing order item.
func TestUpdateOrderItem_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	order := models.Order{Total_amount: 100.00}
	product := models.Product{Name: "Test Product", Price: 10.00}
	db.Create(&order)
	db.Create(&product)
	orderItem := models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 5, Subtotal: 50.00}
	db.Create(&orderItem)

	router.PUT("/orderItems/:id", func(c *gin.Context) {
		UpdateOrderItem(c, db)
	})

	updatedOrderItem := fmt.Sprintf(`{"order_id": %d, "product_id": %d, "quantity": 10, "subtotal": 100.00}`, order.ID, product.ID)
	req, _ := http.NewRequest("PUT", "/orderItems/"+strconv.Itoa(int(orderItem.ID)), bytes.NewBufferString(updatedOrderItem))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.OrderItem
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, 10, response.Quantity)
	assert.Equal(t, 100.00, response.Subtotal)
}

// TestUpdateOrderItem_Invalid checks the error message with invalid data.
func TestUpdateOrderItem_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	order := models.Order{Total_amount: 100.00}
	product := models.Product{Name: "Test Product", Price: 10.00}
	db.Create(&order)
	db.Create(&product)
	orderItem := models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 5, Subtotal: 50.00}
	db.Create(&orderItem)

	router.PUT("/orderItems/:id", func(c *gin.Context) {
		UpdateOrderItem(c, db)
	})

	updatedOrderItem := `{"order_id": "", "product_id": "", "quantity": "", "subtotal": ""}`
	req, _ := http.NewRequest("PUT", "/orderItems/"+strconv.Itoa(int(orderItem.ID)), bytes.NewBufferString(updatedOrderItem))
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

// TestDeleteOrderItem_Valid checks that an order item is deleted from the database.
func TestDeleteOrderItem_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	order := models.Order{Total_amount: 100.00}
	product := models.Product{Name: "Test Product", Price: 10.00}
	db.Create(&order)
	db.Create(&product)
	orderItem := models.OrderItem{Order_ID: uint32(order.ID), Product_ID: uint32(product.ID), Quantity: 5, Subtotal: 50.00}
	db.Create(&orderItem)

	router.DELETE("/orderItems/:id", func(c *gin.Context) {
		DeleteOrderItem(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/orderItems/"+strconv.Itoa(int(orderItem.ID)), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeleteOrderItem_Invalid checks the delete order item with invalid ID.
func TestDeleteOrderItem_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBOrderItem(t)
	defer teardown()

	router.DELETE("/orderItems/:id", func(c *gin.Context) {
		DeleteOrderItem(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/orderItems/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Order item not found")
}
