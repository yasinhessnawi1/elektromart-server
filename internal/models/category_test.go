package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllCategories checks at this function returns all records correctly.
func TestGetAllCategories(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Category A", "Description A").
		AddRow(2, "Category B", "Description B")
	mock.ExpectQuery("^SELECT \\* FROM \"categories\" WHERE").WillReturnRows(rows)

	// Call the function now
	categories, err := GetAllCategories(gormDB)
	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Category A", categories[0].Name)
	assert.Equal(t, "Category B", categories[1].Name)

	// Check all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCategory_SetName check if this function sets the name of category correctly
func TestCategory_SetName(t *testing.T) {
	c := Category{}
	assert.True(t, c.SetName("Valid Category Name"))
	assert.Equal(t, "Valid Category Name", c.Name)
	// Test with invalid name
	assert.False(t, c.SetName(""))
}

// TestCategory_SetDescription check if this function sets the description of category correctly
func TestCategory_SetDescription(t *testing.T) {
	c := Category{}
	assert.True(t, c.SetDescription("Valid Category Description"))
	assert.Equal(t, "Valid Category Description", c.Description)
	// Test with invalid description
	assert.False(t, c.SetDescription(""))
}

// TestCategoryExists checks at this function ensure at a specific category exists.
func TestCategoryExists(t *testing.T) {
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
	mock.ExpectQuery("^SELECT \\* FROM \"categories\"").WithArgs(1, 1).WillReturnRows(rows)

	// Call the function now
	exists := CategoryExists(gormDB, 1)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectQuery("^SELECT \\* FROM \"categories\"").WithArgs(50, 1).WillReturnRows(rows)

	// Call the function now
	exists = CategoryExists(gormDB, 50)
	assert.False(t, exists)
}

// TestSearchCategory checks at this function returns all records correctly.
func TestSearchCategory(t *testing.T) {
	// Create a new instance of sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Search Category", "Matches Criteria")
	mock.ExpectQuery("^SELECT \\* FROM \"categories\" WHERE").WithArgs("%Search Category%").WillReturnRows(rows)

	// Call the function now
	brands, err := SearchCategory(gormDB, map[string]interface{}{"name": "Search Category"})
	assert.NoError(t, err)
	assert.Len(t, brands, 1)
	assert.Equal(t, "Search Category", brands[0].Name)

	// Check all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}
