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

// setupRouterAndDBPayment sets up the router and database in memory, including the migration of Order and Payment models.
// It returns the router, database, and a teardown function to clean up the database after tests finish.
func setupRouterAndDBPayment(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Order{}, &models.Payment{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Order{}, &models.Payment{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

// TestGetPayment_Success checks if GetPayment correctly returns a payment by ID.
// It creates an order and payment in the database, then sends a GET request to retrieve the payment.
// The response is checked for the correct status code and payment details.
func TestGetPayment_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 100.00}
	db.Create(&order)
	payment := models.Payment{Order_ID: uint32(order.ID), Payment_method: "credit card", Amount: 100.00, Payment_date: "2022-01-01", Status: "completed"}
	db.Create(&payment)

	router.GET("/payments/:id", func(c *gin.Context) {
		GetPayment(c, db)
	})

	req, _ := http.NewRequest("GET", "/payments/"+strconv.Itoa(int(payment.ID)), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.Payment
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, payment.ID, response.ID)
}

// TestGetPayment_NotFound checks the response for a non-existent payment.
// It sends a GET request with an invalid ID and checks the response for a 404 status code and an error message.
// The error message should indicate that the payment was not found.
func TestGetPayment_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	router.GET("/payments/:id", func(c *gin.Context) {
		GetPayment(c, db)
	})

	req, _ := http.NewRequest("GET", "/payments/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Payment not found")
}

// TestGetPayments_Success verifies that all payments are correctly retrieved from the database.
// It creates two payments for an order and sends a GET request to retrieve all payments.
// The response is checked for the correct status code and the expected number of payments.
// The test uses the GET /payments route to retrieve all payments.
func TestGetPayments_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 200.00}
	db.Create(&order)
	db.Create(&models.Payment{Order_ID: uint32(order.ID), Payment_method: "credit card", Amount: 100.00, Payment_date: "2022-01-01", Status: "completed"})
	db.Create(&models.Payment{Order_ID: uint32(order.ID), Payment_method: "paypal", Amount: 100.00, Payment_date: "2022-01-02", Status: "completed"})

	router.GET("/payments", func(c *gin.Context) {
		GetPayments(c, db)
	})

	req, _ := http.NewRequest("GET", "/payments", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Payment
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, response, 2)
}

// TestSearchAllPayments_Success checks if payments are correctly retrieved based on search parameters.
// It creates a payment and sends a GET request with a search parameter to retrieve the payment.
// The response is checked for the correct status code and the expected number of payments.
// The test uses the GET /payments/search route to search for payments.
func TestSearchAllPayments_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 300.00}
	db.Create(&order)
	db.Create(&models.Payment{Order_ID: uint32(order.ID), Payment_method: "paypal", Amount: 300.00, Payment_date: "2022-01-01", Status: "completed"})

	router.GET("/payments/search", func(c *gin.Context) {
		SearchAllPayments(c, db)
	})

	req, _ := http.NewRequest("GET", "/payments/search?payment_method=paypal", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response []models.Payment
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, response, 1)
}

// TestCreatePayment_Success checks that a payment can be successfully created with valid data.
// It creates an order in the database and sends a POST request with payment details to create a new payment.
// The response is checked for the correct status code and the expected payment details.
// The test uses the POST /payments route to create a new payment.
func TestCreatePayment_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 500.00}
	db.Create(&order)

	router.POST("/payments", func(c *gin.Context) {
		CreatePayment(c, db)
	})

	newPayment := fmt.Sprintf(`{"order_id": %d, "payment_method": "debit card", "amount": 500.00, "payment_date": "2022-01-01", "status": "completed"}`, order.ID)
	req, _ := http.NewRequest("POST", "/payments", bytes.NewBufferString(newPayment))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.Payment
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, uint32(order.ID), response.Order_ID)
	assert.Equal(t, "completed", response.Status)
}

// TestCreatePayment_InvalidData checks the response when incomplete or incorrect data is sent.
// It sends a POST request with invalid payment details and checks the response for a 400 status code and an error message.
// The error message should indicate that the JSON data is invalid.
// The test uses the POST /payments route to create a new payment.
func TestCreatePayment_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	router.POST("/payments", func(c *gin.Context) {
		CreatePayment(c, db)
	})

	newPayment := `{"order_id": "", "payment_method": "123", "amount": "five hundred", "payment_date": "01-01-2022", "status": "completed"}`
	req, _ := http.NewRequest("POST", "/payments", bytes.NewBufferString(newPayment))
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

// TestUpdatePayment_Valid checks the ability to update an existing payment.
// It creates an order and payment in the database, then sends a PUT request to update the payment status.
// The response is checked for the correct status code and the updated payment status.
// The test uses the PUT /payments/:id route to update an existing payment.
func TestUpdatePayment_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 400.00}
	db.Create(&order)
	payment := models.Payment{Order_ID: uint32(order.ID), Payment_method: "debit card", Amount: 400.00, Payment_date: "2022-01-01", Status: "pending"}
	db.Create(&payment)

	router.PUT("/payments/:id", func(c *gin.Context) {
		UpdatePayment(c, db)
	})

	updatedPayment := fmt.Sprintf(`{"order_id": %d, "payment_method": "debit card", "amount": 400.00, "payment_date": "2022-01-01", "status": "completed"}`, order.ID)
	req, _ := http.NewRequest("PUT", "/payments/"+strconv.Itoa(int(payment.ID)), bytes.NewBufferString(updatedPayment))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.Payment
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "completed", response.Status)
}

// TestUpdatePayment_Invalid checks the error message with invalid data.
// It sends a PUT request with invalid payment details and checks the response for a 400 status code and an error message.
// The error message should indicate that the JSON data is invalid.
// The test uses the PUT /payments/:id route to update an existing payment.
func TestUpdatePayment_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 300.00}
	db.Create(&order)
	payment := models.Payment{Order_ID: uint32(order.ID), Payment_method: "credit card", Amount: 300.00, Payment_date: "2022-01-01", Status: "pending"}
	db.Create(&payment)

	router.PUT("/payments/:id", func(c *gin.Context) {
		UpdatePayment(c, db)
	})

	updatedPayment := `{"order_id": "", "payment_method": "", "amount": "", "payment_date": "", "status": ""}`
	req, _ := http.NewRequest("PUT", "/payments/"+strconv.Itoa(int(payment.ID)), bytes.NewBufferString(updatedPayment))
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

// TestDeletePayment_Valid checks that a payment is deleted from the database.
// It creates an order and payment in the database, then sends a DELETE request to remove the payment.
// The response is checked for the correct status code.
// The test uses the DELETE /payments/:id route to delete an existing payment.
func TestDeletePayment_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	order := models.Order{Total_amount: 100.00}
	db.Create(&order)
	payment := models.Payment{Order_ID: uint32(order.ID), Payment_method: "credit card", Amount: 100.00, Payment_date: "2022-01-01", Status: "completed"}
	db.Create(&payment)

	router.DELETE("/payments/:id", func(c *gin.Context) {
		DeletePayment(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/payments/"+strconv.Itoa(int(payment.ID)), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

// TestDeletePayment_Invalid checks the delete payment with invalid ID.
// It sends a DELETE request with an invalid payment ID and checks the response for a 404 status code and an error message.
// The error message should indicate that the payment was not found.
// The test uses the DELETE /payments/:id route to delete an existing payment.
func TestDeletePayment_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBPayment(t)
	defer teardown()

	router.DELETE("/payments/:id", func(c *gin.Context) {
		DeletePayment(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/payments/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Contains(t, response["error"], "Payment not found")
}
