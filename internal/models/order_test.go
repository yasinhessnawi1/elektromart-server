package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllOrders checks at this function returns all records correctly.
func TestGetAllOrders(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "order_date", "total_amount", "status"}).
		AddRow(1, 1, "2023-01-01", 100.0, "pending").
		AddRow(2, 2, "2023-01-02", 200.0, "completed")
	mock.ExpectQuery("^SELECT \\* FROM \"orders\"").WillReturnRows(rows)

	// Call the function now
	orders, err := GetAllOrders(gormDB)
	assert.NoError(t, err)
	assert.Len(t, orders, 2, "Should fetch two orders")
	assert.Equal(t, uint32(1), orders[0].User_ID, "Check user ID of the first order")
	assert.Equal(t, "2023-01-01", orders[0].Order_date, "Check Order date of the first order")
	assert.Equal(t, float64(100), orders[0].Total_amount, "Check total amount of the first order")
	assert.Equal(t, "pending", orders[0].Status, "Check status of the first order")
}

// TestOrder_SetUserID Checks setting the user ID of the order
func TestOrder_SetUserID(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs(1, 1).
		WillReturnRows(rows)

	// Call the function now
	order := Order{}
	result := order.SetUserID(1, gormDB)
	assert.True(t, result, "User ID should be set when user exists")

	// Not exist user
	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs(50, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = order.SetUserID(50, gormDB)
	assert.False(t, result, "User ID should not be set when user dose not exist")
}

// TestOrder_SetOrderDate checks if this function works correctly.
func TestOrder_SetOrderDate(t *testing.T) {
	order := Order{}
	assert.True(t, order.SetOrderDate("2023-01-01"), "Order date should be valid")
	assert.False(t, order.SetOrderDate("23-01-2023"), "Order date should be invalid")
}

// TestOrder_SetTotalAmount checks if this function works correctly.
func TestOrder_SetTotalAmount(t *testing.T) {
	order := Order{}
	assert.True(t, order.SetTotalAmount(100.0), "Total amount should be valid")
	assert.False(t, order.SetTotalAmount(-1.0), "Total amount should not be set if negative")
}

// TestOrder_SetStatus checks if this function works correctly.
func TestOrder_SetStatus(t *testing.T) {
	order := Order{}
	assert.True(t, order.SetStatus("pending"), "Status date should be valid")
	assert.False(t, order.SetStatus("unknown"), "Status date should be invalid")
}

// TestOrderExists Checks if the function works correctly.
func TestOrderExists(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(1, 1).
		WillReturnRows(rows)

	// Call the function now
	exists := OrderExists(gormDB, 1)
	assert.True(t, exists, "Order should exists")

	// Not exist order
	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(50, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	exists = OrderExists(gormDB, 50)
	assert.False(t, exists, "Order should not exists")
}

// TestSearchOrder Checks if the function works correctly.
func TestSearchOrder(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "order_date", "total_amount", "status"}).
		AddRow(1, 1, "2023-01-01", 100.0, "pending").
		AddRow(2, 2, "2023-01-02", 200.0, "completed")
	mock.ExpectQuery("^SELECT \\* FROM \"orders\"").WithArgs(1).WillReturnRows(rows)

	// Call the function now
	orders, err := SearchOrder(gormDB, map[string]interface{}{"user_id": 1})
	assert.NoError(t, err)
	assert.Len(t, orders, 2, "Should find tow order")
	assert.Equal(t, uint32(1), orders[0].User_ID)

	// Check all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
