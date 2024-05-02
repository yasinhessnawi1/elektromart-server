package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"E-Commerce_Website_Database/internal/models"
)

func setupRouterAndDBProduct(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(&models.Product{}, &models.Brands{}, &models.Category{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Function to clean up the database after tests finish
	teardown := func() {
		if err := db.Migrator().DropTable(&models.Product{}, &models.Brands{}, &models.Category{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}
	return router, db, teardown
}

func TestGetProduct_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	product := models.Product{Name: "Sample Product", Price: 19.99}
	db.Create(&product)

	router.GET("/products/:id", func(c *gin.Context) {
		GetProduct(c, db)
	})

	req, _ := http.NewRequest("GET", fmt.Sprintf("/products/%d", product.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var response models.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, product.Name, response.Name)
}

func TestGetProduct_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	router.GET("/products/:id", func(c *gin.Context) {
		GetProduct(c, db)
	})

	req, _ := http.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetProducts_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	db.Create(&models.Product{Name: "Sample Product 1", Price: 10.00})
	db.Create(&models.Product{Name: "Sample Product 2", Price: 20.00})

	router.GET("/products", func(c *gin.Context) {
		GetProducts(c, db)
	})

	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var products []models.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &products); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Len(t, products, 2)
}

func TestGetProducts_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	router.GET("/products", func(c *gin.Context) {
		GetProducts(c, db)
	})

	req, _ := http.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var products []models.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &products); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, products)
}

func TestSearchProducts_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	// Setups a brand and category as needed
	brand := models.Brands{Name: "Gadget Brand"}
	db.Create(&brand)
	category := models.Category{Name: "Gadgets"}
	db.Create(&category)

	// Create products that should match the search query
	db.Create(&models.Product{Name: "Gadget 1", Description: "Hei", Price: 99.99, Stock_quantity: 50, Brand_ID: uint32(brand.ID), Category_ID: uint32(category.ID)})
	db.Create(&models.Product{Name: "Gadget 2", Description: "Hei", Price: 149.99, Stock_quantity: 100, Brand_ID: uint32(brand.ID), Category_ID: uint32(category.ID)})

	router.GET("/products/search/", func(c *gin.Context) {
		SearchAllProducts(c, db)
	})

	// Define a request with a search query
	req, _ := http.NewRequest("GET", "/products/search/?name=Gadg", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response []models.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Len(t, response, 2) // Expecting two products to match "Hei"
}

func TestSearchProducts_Empty(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	router.GET("/products/search", func(c *gin.Context) {
		c.Request.URL.RawQuery = "name=Nonexistent"
		SearchAllProducts(c, db)
	})

	req, _ := http.NewRequest("GET", "/products/search", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCreateProduct_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	brand := models.Brands{Name: "Tech"}
	db.Create(&brand)
	category := models.Category{Name: "Electronics"}
	db.Create(&category)

	router.POST("/products", func(c *gin.Context) {
		CreateProduct(c, db)
	})

	newProduct := fmt.Sprintf(`{"name": "New Product", "price": 25.50, "description": "A brand new product", "stock_quantity": 100, "brand_id": %d, "category_id": %d}`, brand.ID, category.ID)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(newProduct))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response models.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "New Product", response.Name)
}

func TestCreateProduct_InvalidData(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	brand := models.Brands{Name: "Tech"}
	db.Create(&brand)
	category := models.Category{Name: "Electronics"}
	db.Create(&category)
	router.POST("/products", func(c *gin.Context) {
		CreateProduct(c, db)
	})

	newProduct := fmt.Sprintf(`{"name": "", "price": 25.50, "description": "", "stock_quantity": , "brand_id": %d, "category_id": %d}`, brand.ID, category.ID)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBufferString(newProduct))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateProduct_Success(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	brand := models.Brands{Name: "Gadget Brand"}
	db.Create(&brand)
	category := models.Category{Name: "Gadgets"}
	db.Create(&category)

	product := models.Product{Name: "Old Product", Price: 15.00, Brand_ID: uint32(brand.ID), Category_ID: uint32(category.ID)}
	db.Create(&product)

	router.PUT("/products/:id", func(c *gin.Context) {
		UpdateProduct(c, db)
	})

	updateData := fmt.Sprintf(`{"name": "Updated Product", "price": 20.00,"description": "A brand new product", "stock_quantity": 100, "brand_id": %d, "category_id": %d}`, uint32(category.ID), uint32(brand.ID))
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%d", product.ID), bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	assert.Equal(t, "Updated Product", response.Name)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	router.PUT("/products/:id", func(c *gin.Context) {
		UpdateProduct(c, db)
	})

	updateData := `{"name": "Updated Product", "price": 50.00}`
	req, _ := http.NewRequest("PUT", "/products/999", bytes.NewBufferString(updateData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteProduct_Valid(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	product := models.Product{Name: "Delete Product", Price: 30.00}
	db.Create(&product)

	router.DELETE("/products/:id", func(c *gin.Context) {
		DeleteProduct(c, db)
	})

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%d", product.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteProduct_Invalid(t *testing.T) {
	router, db, teardown := setupRouterAndDBProduct(t)
	defer teardown()

	router.DELETE("/products/:id", func(c *gin.Context) {
		DeleteProduct(c, db)
	})

	req, _ := http.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
