package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllPayments ensures that all payments are retrieved correctly.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestGetAllPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "order_id", "payment_method", "amount", "payment_date", "status"}).
		AddRow(1, 1, "Credit Card", 100.0, "2021-04-21", "Completed").
		AddRow(2, 2, "PayPal", 200.0, "2021-04-22", "Pending")
	mock.ExpectQuery("^SELECT \\* FROM \"payments\"").WillReturnRows(rows)

	payments, err := GetAllPayments(gormDB)
	assert.NoError(t, err)
	assert.Len(t, payments, 2, "Should fetch two payments")
}

// TestPayment_SetOrderID checks that the order ID is set only after validating the order exists.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestPayment_SetOrderID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	payment := Payment{}
	result := payment.SetOrderID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"orders\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = payment.SetOrderID(99, gormDB)
	assert.False(t, result)
}

// TestPayment_SetPaymentMethod checks that the payment method is set only if it is valid.
// It creates a new instance of the Payment struct and calls the SetPaymentMethod function with a valid method.
// It then checks if the method was set correctly and if the function returned true.
// It repeats the process with an invalid method and checks if the method was not set and the function returned false.
func TestPayment_SetPaymentMethod(t *testing.T) {
	payment := Payment{}
	assert.False(t, payment.SetPaymentMethod("Invalid Method"))
	assert.True(t, payment.SetPaymentMethod("Credit Card"))
}

// TestPayment_SetAmount checks the validity and setting of the payment amount.
// It creates a new instance of the Payment struct and calls the SetAmount function with a valid amount.
// It then checks if the amount was set correctly and if the function returned true.
// It repeats the process with an invalid amount and checks if the amount was not set and the function returned false.
func TestPayment_SetAmount(t *testing.T) {
	payment := Payment{}
	assert.False(t, payment.SetAmount(-100.0))
	assert.True(t, payment.SetAmount(100.0))
}

// TestPayment_SetPaymentDate checks the date setting and validation logic.
// It creates a new instance of the Payment struct and calls the SetPaymentDate function with a valid date.
// It then checks if the date was set correctly and if the function returned true.
// It repeats the process with an invalid date and checks if the date was not set and the function returned false.
func TestPayment_SetPaymentDate(t *testing.T) {
	payment := Payment{}
	assert.False(t, payment.SetPaymentDate("not-a-date"))
	assert.True(t, payment.SetPaymentDate("2021-04-21"))
}

// TestPayment_SetStatus checks the status setting and validation logic.
// It creates a new instance of the Payment struct and calls the SetStatus function with a valid status.
// It then checks if the status was set correctly and if the function returned true.
// It repeats the process with an invalid status and checks if the status was not set and the function returned false.
func TestPayment_SetStatus(t *testing.T) {
	payment := Payment{}
	assert.False(t, payment.SetStatus("unknown"), "Status should be invalid")
	assert.True(t, payment.SetStatus("completed"), "Expected 'completed' to be a valid status")

}

// TestPaymentExists checks if a payment exists in the database by its ID.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestPaymentExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"payments\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	assert.True(t, PaymentExists(gormDB, 1))

	mock.ExpectQuery("^SELECT \\* FROM \"payments\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	assert.False(t, PaymentExists(gormDB, 2))
}

// TestSearchPayment checks the query functionality to search for payments based on provided parameters.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestSearchPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "order_id", "payment_method", "amount", "payment_date", "status"}).
		AddRow(1, 1, "Credit Card", 100.0, "2021-04-21", "Completed")
	mock.ExpectQuery("^SELECT \\* FROM \"payments\" WHERE").
		WithArgs("%credit card%").
		WillReturnRows(rows)

	searchParams := map[string]interface{}{"payment_method": "credit card"}
	payments, err := SearchPayment(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, payments, 1, "Should find one payment")
}
