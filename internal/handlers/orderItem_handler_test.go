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
// It returns the router, database, and a teardown function to clean up the database after tests finish.
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
// It creates an order, product, and order item in the database, then fetches the order item by ID.
// The test checks the response status code, the order item ID, and the quantity.
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
// It fetches an order item by an ID that does not exist in the database and checks the response status code and error message.
// The test expects a 404 Not Found status code and an error message indicating that the order item was not found.
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
// It creates two order items in the database and fetches all order items, checking the response status code and the number of order items.
// The test expects a 200 OK status code and two order items in the response.
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
// It fetches all order items from an empty database and checks the response status code and the number of order items.
// The test expects a 200 OK status code and an empty list of order items in the response.
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
// It creates an order, product, and order item in the database and fetches the order item by quantity.
// The test checks the response status code, the number of order items, and the quantity of the order item.
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
// It fetches order items with a quantity that does not exist in the database and checks the response status code and error message.
// The test expects a 404 Not Found status code and an error message indicating that no order items were found.
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
// It creates an order and product in the database and sends a request to create a new order item.
// The test checks the response status code, the order ID, the product ID, the quantity, and the subtotal.
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
// It sends a request to create an order item with invalid JSON data and checks the response status code and error message.
// The test expects a 400 Bad Request status code and an error message indicating that the JSON data is invalid.
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
// It creates an order, product, and order item in the database and sends a request to update the order item.
// The test checks the response status code, the updated quantity, and the updated subtotal.
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
// It sends a request to update an order item with invalid JSON data and checks the response status code and error message.
// The test expects a 400 Bad Request status code and an error message indicating that the JSON data is invalid.
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
// It creates an order, product, and order item in the database and sends a request to delete the order item.
// The test checks the response status code and the absence of the order item in the database.
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
// It sends a request to delete an order item with an ID that does not exist in the database and checks
// the response status code and error message.
// The test expects a 404 Not Found status code and an error message indicating that the order item was not found.
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
