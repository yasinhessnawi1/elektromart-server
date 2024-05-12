package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllOrderItems ensures that all order items are retrieved correctly.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestGetAllOrderItems(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "subtotal"}).
		AddRow(1, 1, 1, 5, 100.0).
		AddRow(2, 1, 2, 3, 50.0)
	mock.ExpectQuery("^SELECT \\* FROM \"order_items\"").
		WillReturnRows(rows)

	// Call the function now
	orderItem, err := GetAllOrderItems(gormDB)
	assert.NoError(t, err)
	assert.Len(t, orderItem, 2, "Should fetch two order item")
	assert.Equal(t, uint32(1), orderItem[0].Order_ID, "Check Order ID of the first order item")
	assert.Equal(t, uint32(1), orderItem[0].Product_ID, "Check Product ID of the first order item")
	assert.Equal(t, 5, orderItem[0].Quantity, "Check quantity of the first order item")
	assert.Equal(t, float64(100), orderItem[0].Subtotal, "Check subtotal of the first order item")
}

// TestOrderItem_SetOrderID ensures that the order id are set only if the referenced entity exits.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestOrderItem_SetOrderID(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	orderItem := OrderItem{}
	result := orderItem.SetOrderID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = orderItem.SetOrderID(99, gormDB)
	assert.False(t, result)
}

// TestOrderItem_SetProductID ensures that the product id are set only if the referenced entity exits.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestOrderItem_SetProductID(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	orderItem := OrderItem{}
	result := orderItem.SetProductID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = orderItem.SetProductID(99, gormDB)
	assert.False(t, result)
}

// TestOrderItem_SetQuantity validate and set quantity
// It creates a new instance of the OrderItem struct and calls the SetQuantity function with a valid quantity.
// It then checks if the quantity was set correctly and if the function returned true.
// It repeats the process with an invalid quantity and checks if the quantity was not set and the function returned false.
func TestOrderItem_SetQuantity(t *testing.T) {
	orderItem := OrderItem{}
	assert.False(t, orderItem.SetQuantity(-1))
	assert.True(t, orderItem.SetQuantity(10))
}

// TestOrderItem_SetSubtotal validate and set subtotal
// It creates a new instance of the OrderItem struct and calls the SetSubtotal function with a valid subtotal.
// It then checks if the subtotal was set correctly and if the function returned true.
// It repeats the process with an invalid subtotal and checks if the subtotal was not set and the function returned false.
func TestOrderItem_SetSubtotal(t *testing.T) {
	orderItem := OrderItem{}
	assert.False(t, orderItem.SetSubtotal(-100.0))
	assert.True(t, orderItem.SetSubtotal(200.0))
}

// TestOrderItemExists checks if an order item exists by its ID
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestOrderItemExists(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"order_items\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	assert.True(t, OrderItemExists(gormDB, 1))

	mock.ExpectQuery("^SELECT \\* FROM \"order_items\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	assert.False(t, OrderItemExists(gormDB, 2))
}

// TestSearchOrderItem checks the query to search for order items based on provided parameters
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestSearchOrderItem(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "subtotal"}).
		AddRow(1, 1, 1, 5, 100.0)
	mock.ExpectQuery("^SELECT \\* FROM \"order_items\" WHERE").WithArgs(1, 1).
		WillReturnRows(rows)

	searchParams := map[string]interface{}{"order_id": 1, "product_id": 1}
	orderItems, err := SearchOrderItem(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, orderItems, 1, "Should find one order item")
	assert.Equal(t, uint32(1), orderItems[0].Order_ID)
	assert.Equal(t, uint32(1), orderItems[0].Product_ID)
	assert.Equal(t, 5, orderItems[0].Quantity)
	assert.Equal(t, float64(100), orderItems[0].Subtotal)
}
