package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllBrands checks at this function returns all records correctly.
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
func TestBrands_SetName(t *testing.T) {
	b := Brands{}
	assert.True(t, b.SetName("Valid Brand Name"))
	assert.Equal(t, "Valid Brand Name", b.Name)
	// Test wth invalid name
	assert.False(t, b.SetName(""))
}

// TestBrands_SetDescription check if this function sets the description of brand correctly
func TestBrands_SetDescription(t *testing.T) {
	b := Brands{}
	assert.True(t, b.SetDescription("Valid Brand Description"))
	assert.Equal(t, "Valid Brand Description", b.Description)
	// Test wth invalid description
	assert.False(t, b.SetDescription(""))
}

// TestBrandExists checks at this function ensure at a specific brand exists.
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
	mock.ExpectQuery("^SELECT \\* FROM \"brands\" WHERE").WithArgs("Search Brand", 1).WillReturnRows(rows)

	// Call the function now
	brands, err := SearchBrand(gormDB, map[string]interface{}{"name": "Search Brand"})
	assert.NoError(t, err)
	assert.NotNil(t, brands)
	assert.Equal(t, "Search Brand", brands.Name)

	// Check all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
