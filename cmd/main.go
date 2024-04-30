package main

import (
	"E-Commerce_Website_Database/internal/config"
	"E-Commerce_Website_Database/internal/handlers"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func main() {
	config.LoadConfig()
	port := config.GetConfig("PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	r := gin.Default()

	// Configuring CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true                                                              // Allow all origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"} // Allow all methods
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	corsConfig.AddExposeHeaders("Access-Control-Allow-Origin") // Add this line
	// Allow headers
	r.Use(LoggerMiddleware())
	r.Use(cors.New(corsConfig))
	setupRoutes(r, db)
	r.POST("/login", func(context *gin.Context) { handlers.PostLogin(context, db) })
	r.GET("/protected", tools.TokenAuthMiddleware(), func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{"username": username, "message": "Welcome to the protected route!"})
	})
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, db *gorm.DB) {
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
	router.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
	router.GET("/users/:id", func(c *gin.Context) { handlers.GetUser(c, db) })
	router.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
	router.PUT("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, db) })
	router.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, db) })
	// Here you should use Query Param Like :search-users/?username={The username}  or search-users/?email={The email}
	//`or by first name , last name , or address`.
	router.GET("/search-users/", func(c *gin.Context) { handlers.SearchAllUsers(c, db) })

	router.GET("/products", func(c *gin.Context) { handlers.GetProducts(c, db) })
	router.GET("/products/:id", func(c *gin.Context) { handlers.GetProduct(c, db) })
	router.POST("/products", func(c *gin.Context) { handlers.CreateProduct(c, db) })
	router.PUT("/products/:id", func(c *gin.Context) { handlers.UpdateProduct(c, db) })
	router.DELETE("/products/:id", func(c *gin.Context) { handlers.DeleteProduct(c, db) })
	// Here you should use Query Param Like :search-products/?name={The name of product}  or search-users/?price={The price}
	//`or by brand id , category id`.
	router.GET("/search-products/:any", func(c *gin.Context) { handlers.SearchAllProducts(c, db) })

	router.GET("/brand", func(c *gin.Context) { handlers.GetBrands(c, db) })
	router.GET("/brand/:id", func(c *gin.Context) { handlers.GetBrand(c, db) })
	router.POST("/brand", func(c *gin.Context) { handlers.CreateBrand(c, db) })
	router.PUT("/brand/:id", func(c *gin.Context) { handlers.UpdateBrand(c, db) })
	router.DELETE("/brand/:id", func(c *gin.Context) { handlers.DeleteBrand(c, db) })

	router.GET("/categories", func(c *gin.Context) { handlers.GetCategories(c, db) })
	router.GET("/categories/:id", func(c *gin.Context) { handlers.GetCategory(c, db) })
	router.POST("/categories", func(c *gin.Context) { handlers.CreateCategory(c, db) })
	router.PUT("/categories/:id", func(c *gin.Context) { handlers.UpdateCategory(c, db) })
	router.DELETE("/categories/:id", func(c *gin.Context) { handlers.DeleteCategory(c, db) })

	router.GET("/orders", func(c *gin.Context) { handlers.GetOrders(c, db) })
	router.GET("/orders/:id", func(c *gin.Context) { handlers.GetOrder(c, db) })
	router.POST("/orders", func(c *gin.Context) { handlers.CreateOrder(c, db) })
	router.PUT("/orders/:id", func(c *gin.Context) { handlers.UpdateOrder(c, db) })
	router.DELETE("/orders/:id", func(c *gin.Context) { handlers.DeleteOrder(c, db) })

	router.GET("/orderItems", func(c *gin.Context) { handlers.GetOrderItems(c, db) })
	router.GET("/orderItems/:id", func(c *gin.Context) { handlers.GetOrderItem(c, db) })
	router.POST("/orderItems", func(c *gin.Context) { handlers.CreateOrderItem(c, db) })
	router.PUT("/orderItems/:id", func(c *gin.Context) { handlers.UpdateOrderItem(c, db) })
	router.DELETE("/orderItems/:id", func(c *gin.Context) { handlers.DeleteOrderItem(c, db) })

	router.GET("/payments", func(c *gin.Context) { handlers.GetPayments(c, db) })
	router.GET("/payments/:id", func(c *gin.Context) { handlers.GetPayment(c, db) })
	router.POST("/payments", func(c *gin.Context) { handlers.CreatePayment(c, db) })
	router.PUT("/payments/:id", func(c *gin.Context) { handlers.UpdatePayment(c, db) })
	router.DELETE("/payments/:id", func(c *gin.Context) { handlers.DeletePayment(c, db) })
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Request received")
		c.Next()
		fmt.Println("Headers:", c.Writer.Header())
	}
}
