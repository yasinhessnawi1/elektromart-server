package handlers

/*
import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router = SetupRoutes()

func SetupRoutes() *gin.Engine {
	dsn := "root:@tcp(localhost:8000)/eCommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Create a new router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to ElectroMart API"})
	})
	router.HandleMethodNotAllowed = true

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found or method not allowed"})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	})

	// Setup routes in /cmd/main.go or wherever you configure routes
	router.GET("/users", func(c *gin.Context) { GetUsers(c, db) })
	router.GET("/users/:id", func(c *gin.Context) { GetUser(c, db) })
	router.POST("/users", func(c *gin.Context) { CreateUser(c, db) })
	router.PUT("/users/:id", func(c *gin.Context) { UpdateUser(c, db) })
	router.DELETE("/users/:id", func(c *gin.Context) { DeleteUser(c, db) })
	// Here you should use Query Param Like :search-users/?username={The username}  or search-users/?email={The email}
	//`or by first name , last name , or address`.
	router.GET("/search-users/", func(c *gin.Context) { SearchAllUsers(c, db) })

	router.GET("/products", func(c *gin.Context) { GetProducts(c, db) })
	router.GET("/products/:id", func(c *gin.Context) { GetProduct(c, db) })
	router.POST("/products", func(c *gin.Context) { CreateProduct(c, db) })
	router.PUT("/products/:id", func(c *gin.Context) { UpdateProduct(c, db) })
	router.DELETE("/products/:id", func(c *gin.Context) { DeleteProduct(c, db) })
	// Here you should use Query Param Like :search-products/?name={The name of product}  or search-users/?price={The price}
	//`or by brand id , category id`.
	router.GET("/search-products/", func(c *gin.Context) { SearchAllProducts(c, db) })

	router.GET("/brand", func(c *gin.Context) { GetBrands(c, db) })
	router.GET("/brand/:id", func(c *gin.Context) { GetBrand(c, db) })
	router.POST("/brand", func(c *gin.Context) { CreateBrand(c, db) })
	router.PUT("/brand/:id", func(c *gin.Context) { UpdateBrand(c, db) })
	router.DELETE("/brand/:id", func(c *gin.Context) { DeleteBrand(c, db) })

	router.GET("/categories", func(c *gin.Context) { GetCategories(c, db) })
	router.GET("/categories/:id", func(c *gin.Context) { GetCategory(c, db) })
	router.POST("/categories", func(c *gin.Context) { CreateCategory(c, db) })
	router.PUT("/categories/:id", func(c *gin.Context) { UpdateCategory(c, db) })
	router.DELETE("/categories/:id", func(c *gin.Context) { DeleteCategory(c, db) })

	router.GET("/orders", func(c *gin.Context) { GetOrders(c, db) })
	router.GET("/orders/:id", func(c *gin.Context) { GetOrder(c, db) })
	router.POST("/orders", func(c *gin.Context) { CreateOrder(c, db) })
	router.PUT("/orders/:id", func(c *gin.Context) { UpdateOrder(c, db) })
	router.DELETE("/orders/:id", func(c *gin.Context) { DeleteOrder(c, db) })

	router.GET("/orderItems", func(c *gin.Context) { GetOrderItems(c, db) })
	router.GET("/orderItems/:id", func(c *gin.Context) { GetOrderItem(c, db) })
	router.POST("/orderItems", func(c *gin.Context) { CreateOrderItem(c, db) })
	router.PUT("/orderItems/:id", func(c *gin.Context) { UpdateOrderItem(c, db) })
	router.DELETE("/orderItems/:id", func(c *gin.Context) { DeleteOrderItem(c, db) })

	router.GET("/payments", func(c *gin.Context) { GetPayments(c, db) })
	router.GET("/payments/:id", func(c *gin.Context) { GetPayment(c, db) })
	router.POST("/payments", func(c *gin.Context) { CreatePayment(c, db) })
	router.PUT("/payments/:id", func(c *gin.Context) { UpdatePayment(c, db) })
	router.DELETE("/payments/:id", func(c *gin.Context) { DeletePayment(c, db) })
	return router
}

func TestGetBrandWithGin(t *testing.T) {
	t.Run("Test GetBrand with valid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brand/1524773729", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected HTTP status 200 OK, got %d", w.Code)
		}
	})

	t.Run("Test GetBrand with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brand/148683", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected HTTP status 404 with an error, got %d", w.Code)
		}
	})
}

func TestGetBrandsWithGin(t *testing.T) {

	t.Run("Test GetBrands with valid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brand/1524773729", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected HTTP status 200 OK, got %d", w.Code)
		} else {
			t.Log("Test passed")
		}
	})

	t.Run("Test GetBrands with invalid ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brand/148683", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected HTTP status 404 with an error, got %d", w.Code)
		} else {
			t.Log("Test passed")
		}

	})
}

func TestCreateBrand(t *testing.T) {
	t.Run("Test CreateBrand with valid data", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Replace nil with the correct data
		body := map[string]interface{}{
			"name":        "Test Brand",
			"description": "This is a test brand",
		}
		correctData, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(correctData))
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("Expected HTTP status 201 Created, got %d", w.Code)
		} else {
			t.Log("Test passed")
		}
	})

	t.Run("Test CreateBrand with invalid data", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := map[string]interface{}{
			"test": "test",
		}
		incorrectData, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(incorrectData))
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", w.Code)
		} else {
			t.Log("Test passed")
		}
	})
}

func TestUpdateBrand(t *testing.T) {
	t.Run("Test UpdateBrand with valid data", func(t *testing.T) {
		w := httptest.NewRecorder()

		Body := map[string]interface{}{
			"name":        "Test Brand",
			"description": "This is a test brand",
		}
		bodyData, _ := json.Marshal(Body)
		req, _ := http.NewRequest("PUT", "/brand/1524773729", bytes.NewBuffer(bodyData))
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected HTTP status 200 OK, got %d", w.Code)
		} else {
			t.Log("Test passed")
		}
	})

	t.Run("Test UpdateBrand with invalid data", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/brand/148683", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected HTTP status 500 Internal Server Error, got %d", w.Code)
		} else {
			t.Log("Test passed")
		}
	})
}

*/
