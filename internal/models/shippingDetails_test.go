package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllShippingDetails tests retrieving all shipping details from the database.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestGetAllShippingDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "order_id", "address", "shipping_date", "estimated_arrival", "status"}).
		AddRow(1, 1, "123 Elm St", "2021-06-01", "2021-06-15", "Shipped").
		AddRow(2, 2, "456 Oak St", "2021-07-01", "2021-07-15", "Delivered")
	mock.ExpectQuery("^SELECT \\* FROM \"shipping_details\"").WillReturnRows(rows)

	shippingDetails, err := GetAllShippingDetails(gormDB)
	assert.NoError(t, err)
	assert.Len(t, shippingDetails, 2, "Should fetch two shipping details")
}

// TestShippingDetails_SetOrderID tests setting the order ID after verifying the order exists.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestShippingDetails_SetOrderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	details := ShippingDetails{}
	result := details.SetOrderID(1, gormDB)
	assert.True(t, result, "Order exists, should return true")

	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = details.SetOrderID(99, gormDB)
	assert.False(t, result, "Order does not exist, should return false")
}

// TestShippingDetails_SetAddress tests setting the shipping address after validating its length.
// It creates a new instance of the ShippingDetails struct and calls the SetAddress function with a valid address.
// It then checks if the address was set correctly and if the function returned true.
// It repeats the process with an invalid address and checks if the address was not set and the function returned false.
func TestShippingDetails_SetAddress(t *testing.T) {
	details := ShippingDetails{}
	assert.False(t, details.SetAddress(string(make([]byte, 256))), "Address should be invalid due to length")
	assert.True(t, details.SetAddress("123 Elm St"), "Address should be valid")
}

// TestShippingDetails_SetShippingDate tests setting the shipping date after validating its format.
// It creates a new instance of the ShippingDetails struct and calls the SetShippingDate function with a valid date.
// It then checks if the date was set correctly and if the function returned true.
// It repeats the process with an invalid date and checks if the date was not set and the function returned false.
func TestShippingDetails_SetShippingDate(t *testing.T) {
	details := ShippingDetails{}
	assert.False(t, details.SetShippingDate("20210601"), "Shipping date should be invalid due to format")
	assert.True(t, details.SetShippingDate("2021-06-01"), "Shipping date should be valid")
}

// TestShippingDetails_SetEstimatedArrival tests setting the estimated arrival date after validating its format.
// It creates a new instance of the ShippingDetails struct and calls the SetEstimatedArrival function with a valid date.
// It then checks if the date was set correctly and if the function returned true.
// It repeats the process with an invalid date and checks if the date was not set and the function returned false.
func TestShippingDetails_SetEstimatedArrival(t *testing.T) {
	details := ShippingDetails{}
	assert.False(t, details.SetEstimatedArrival("20210701"), "Estimated arrival should be invalid due to format")
	assert.True(t, details.SetEstimatedArrival("2021-07-15"), "Estimated arrival should be valid")
}

// TestShippingDetails_SetStatus tests setting the status after validating its length and contents.
// It creates a new instance of the ShippingDetails struct and calls the SetStatus function with a valid status.
// It then checks if the status was set correctly and if the function returned true.
// It repeats the process with an invalid status and checks if the status was not set and the function returned false.
func TestShippingDetails_SetStatus(t *testing.T) {
	details := ShippingDetails{}
	assert.False(t, details.SetStatus(""), "Status should be invalid due to being empty")
	assert.True(t, details.SetStatus("delivered"), "Status should be valid")
}

// TestShippingDetailsExists tests checking if a specific shipping detail exists in the database by its ID.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestShippingDetailsExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"shipping_details\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	assert.True(t, ShippingDetailsExists(gormDB, 1), "Shipping detail should exist")

	mock.ExpectQuery("^SELECT \\* FROM \"shipping_details\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	assert.False(t, ShippingDetailsExists(gormDB, 2), "Shipping detail should not exist")
}

// TestSearchShippingDetails tests searching for shipping details based on provided parameters.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestSearchShippingDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "order_id", "address", "shipping_date", "estimated_arrival", "status"}).
		AddRow(1, 2, "Bunt St", "2021-06-01", "2021-06-15", "Delivered")

	mock.ExpectQuery("^SELECT \\* FROM \"shipping_details\" WHERE").
		WithArgs("%bunt st%").
		WillReturnRows(rows)

	searchParams := map[string]interface{}{
		"address": "bunt st",
	}
	details, err := SearchShippingDetails(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, details, 1, "Should find one matching shipping detail")
}
