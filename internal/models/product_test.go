package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllProducts tests that all products are retrieved from the database correctly.
func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock_quantity", "brand_id", "category_id"}).
		AddRow(1, "Product 1", "Description 1", 10.50, 5, 1, 1).
		AddRow(2, "Product 2", "Description 2", 20.75, 8, 1, 2)
	mock.ExpectQuery("^SELECT \\* FROM \"products\"").WillReturnRows(rows)

	products, err := GetAllProducts(gormDB)
	assert.NoError(t, err)
	assert.Len(t, products, 2, "Should fetch two products")
}

// TestProduct_SetName tests setting a product's name after validating its length.
func TestProduct_SetName(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetName(string(make([]byte, 256))), "Name should be invalid due to length")
	assert.True(t, product.SetName("Valid Name"), "Name should be valid")
}

// TestProduct_SetDescription tests setting a product's description after validating its length.
func TestProduct_SetDescription(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetDescription(string(make([]byte, 1001))), "Description should be invalid due to length")
	assert.True(t, product.SetDescription("Valid Description"), "Description should be valid")
}

// TestProduct_SetPrice tests setting a product's price after validating it as a positive float.
func TestProduct_SetPrice(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetPrice(-10.0), "Price should be invalid because it is negative")
	assert.True(t, product.SetPrice(100.0), "Price should be valid")
}

// TestProduct_SetStockQuantity tests setting a product's stock quantity after validating it as a non-negative integer.
func TestProduct_SetStockQuantity(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetStockQuantity(-1), "Stock quantity should be invalid because it is negative")
	assert.True(t, product.SetStockQuantity(50), "Stock quantity should be valid")
}

// TestProduct_SetBrandID tests setting a product's brand ID after verifying the existence of the brand.
func TestProduct_SetBrandID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"brands\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	product := Product{}
	result := product.SetBrandID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"brands\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = product.SetBrandID(2, gormDB)
	assert.False(t, result)
}

// TestProduct_SetCategoryID tests setting a product's category ID after verifying the existence of the category.
func TestProduct_SetCategoryID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"categories\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	product := Product{}
	result := product.SetCategoryID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"categories\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = product.SetCategoryID(2, gormDB)
	assert.False(t, result)
}

// TestProductExists tests checking if a specific product exists in the database by its ID.
func TestProductExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	assert.True(t, ProductExists(gormDB, 1), "Product should exist")

	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	assert.False(t, ProductExists(gormDB, 2), "Product should not exist")
}

// TestSearchProduct tests the functionality to search for products based on provided parameters.
func TestSearchProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock_quantity", "brand_id", "category_id"}).
		AddRow(1, "Searchable Product", "Description", 10.0, 5, 1, 1)
	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").
		WithArgs("%searchable%", 10.0, 5).
		WillReturnRows(rows)

	searchParams := map[string]interface{}{
		"name":           "searchable",
		"price":          10.0,
		"stock_quantity": 5,
	}
	products, err := SearchProduct(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, products, 1, "Should find one matching product")
}
