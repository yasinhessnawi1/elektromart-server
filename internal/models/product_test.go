package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllProducts tests that all products are retrieved from the database correctly.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
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
// It creates a new instance of the Product struct and calls the SetName function with a valid name.
// It then checks if the name was set correctly and if the function returned true.
// It repeats the process with an invalid name and checks if the name was not set and the function returned false.
func TestProduct_SetName(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetName(string(make([]byte, 256))), "Name should be invalid due to length")
	assert.True(t, product.SetName("Valid Name"), "Name should be valid")
}

// TestProduct_SetDescription tests setting a product's description after validating its length.
// It creates a new instance of the Product struct and calls the SetDescription function with a valid description.
// It then checks if the description was set correctly and if the function returned true.
// It repeats the process with an invalid description and checks if the description was not set and the function returned false.
func TestProduct_SetDescription(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetDescription(string(make([]byte, 1001))), "Description should be invalid due to length")
	assert.True(t, product.SetDescription("Valid Description"), "Description should be valid")
}

// TestProduct_SetPrice tests setting a product's price after validating it as a positive float.
// It creates a new instance of the Product struct and calls the SetPrice function with a valid price.
// It then checks if the price was set correctly and if the function returned true.
// It repeats the process with an invalid price and checks if the price was not set and the function returned false.
func TestProduct_SetPrice(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetPrice(-10.0), "Price should be invalid because it is negative")
	assert.True(t, product.SetPrice(100.0), "Price should be valid")
}

// TestProduct_SetStockQuantity tests setting a product's stock quantity after validating it as a non-negative integer.
// It creates a new instance of the Product struct and calls the SetStockQuantity function with a valid quantity.
// It then checks if the quantity was set correctly and if the function returned true.
// It repeats the process with an invalid quantity and checks if the quantity was not set and the function returned false.
func TestProduct_SetStockQuantity(t *testing.T) {
	product := Product{}
	assert.False(t, product.SetStockQuantity(-1), "Stock quantity should be invalid because it is negative")
	assert.True(t, product.SetStockQuantity(50), "Stock quantity should be valid")
}

// TestProduct_SetBrandID tests setting a product's brand ID after verifying the existence of the brand.
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
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
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
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
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
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
// It creates a new instance of sql mock and sets up expectations for the query.
// It then calls the function and checks if the returned data matches the expected data.
// Finally, it checks if all the expectations were met.
func TestSearchProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "stock_quantity", "brand_id", "category_id"}).
		AddRow(1, "Searchable Product", "Description", 10.0, 5, 1, 1)

	// match the actual query that is being executed
	mock.ExpectQuery(`^SELECT "products"."id","products"."created_at","products"."updated_at","products"."deleted_at","products"."name","products"."description","products"."price","products"."stock_quantity","products"."brand_id","products"."category_id" FROM "products" JOIN brands ON brands.id = products.brand_id JOIN categories ON categories.id = products.category_id WHERE products.name LIKE \$1 AND "products"."deleted_at" IS NULL$`).
		WithArgs("%searchable%").
		WillReturnRows(rows)

	searchParams := map[string]interface{}{
		"name": "searchable",
	}
	products, err := SearchProduct(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, products, 1, "Should find one matching product")
}
