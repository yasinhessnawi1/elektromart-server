package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllBrands checks at this function returns all records correctly.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestGetAllBrands(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Brand A", "Description A").
		AddRow(2, "Brand B", "Description B")
	mock.ExpectQuery("^SELECT \\* FROM \"brands\" WHERE \"brands\".\"deleted_at\" IS NULL").WillReturnRows(rows)

	// Call the function now
	brands, err := GetAllBrands(gormDB)
	assert.NoError(t, err)
	assert.Len(t, brands, 2)
	assert.Equal(t, "Brand A", brands[0].Name)
	assert.Equal(t, "Brand B", brands[1].Name)

	// Check all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestBrands_SetName check if this function sets the name of brand correctly
// It creates a new instance of the Brands struct and calls the SetName function with a valid name.
// It then checks if the name was set correctly and if the function returned true.
// It repeats the process with an invalid name and checks if the name was not set and the function returned false.
func TestBrands_SetName(t *testing.T) {
	b := Brands{}
	assert.True(t, b.SetName("Valid Brand Name"))
	assert.Equal(t, "Valid Brand Name", b.Name)
	// Test wth invalid name
	assert.False(t, b.SetName(""))
}

// TestBrands_SetDescription check if this function sets the description of brand correctly
// It creates a new instance of the Brands struct and calls the SetDescription function with a valid description.
// It then checks if the description was set correctly and if the function returned true.
// It repeats the process with an invalid description and checks if the description was not set and the function returned false.
func TestBrands_SetDescription(t *testing.T) {
	b := Brands{}
	assert.True(t, b.SetDescription("Valid Brand Description"))
	assert.Equal(t, "Valid Brand Description", b.Description)
	// Test wth invalid description
	assert.False(t, b.SetDescription(""))
}

// TestBrandExists checks at this function ensure at a specific brand exists.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestBrandExists(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery("^SELECT \\* FROM \"brands\"").WithArgs(1, 1).WillReturnRows(rows)

	// Call the function now
	exists := BrandExists(gormDB, 1)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectQuery("^SELECT \\* FROM \"brands\"").WithArgs(50, 1).WillReturnRows(rows)

	// Call the function now
	exists = BrandExists(gormDB, 50)
	assert.False(t, exists)
}

// TestSearchBrand checks at this function returns all records correctly.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestSearchBrand(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Search Brand", "Matches Criteria")
	mock.ExpectQuery("^SELECT \\* FROM \"brands\" WHERE").WithArgs("Search Brand").WillReturnRows(rows)

	// Call the function now
	brands, err := SearchBrand(gormDB, map[string]interface{}{"name": "Search Brand"})
	assert.NoError(t, err)
	assert.Len(t, brands, 1)
	assert.Equal(t, "Search Brand", brands.Name)

	// Check all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
