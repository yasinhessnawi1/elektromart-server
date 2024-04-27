package main

import (
	"E-Commerce_Website_Database/internal/config"
	"E-Commerce_Website_Database/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // Automatically loads the .env file
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// main initializes the application, sets up database connections, and starts the HTTP server.
func main() {
	config.LoadConfig() // Load environment variables from .env file

	// Set up database connection using the MySQL driver.
	port := config.GetConfig("DATABASE_PORT")
	dsn := "root:@tcp(localhost:" + port + ")/eCommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up the Gin router and middleware.
	router := gin.Default()
	router.Use(CORSMiddleware()) // Apply CORS middleware to allow cross-origin requests.
	setupRoutes(router, db)      // Set up API routes.
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// setupRoutes defines all the routes and their handlers for the application.
func setupRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to ElectroMart API"})
	})
	// Handle requests for non-existent routes.
	router.HandleMethodNotAllowed = true
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "endpoint not found or method not allowed"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	})

	// User routes
	router.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
	router.GET("/users/:id", func(c *gin.Context) { handlers.GetUser(c, db) })
	router.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
	router.PUT("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, db) })
	router.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, db) })
	router.GET("/search-users/", func(c *gin.Context) { handlers.SearchAllUsers(c, db) })

	// Product routes
	router.GET("/products", func(c *gin.Context) { handlers.GetProducts(c, db) })
	router.GET("/products/:id", func(c *gin.Context) { handlers.GetProduct(c, db) })
	router.POST("/products", func(c *gin.Context) { handlers.CreateProduct(c, db) })
	router.PUT("/products/:id", func(c *gin.Context) { handlers.UpdateProduct(c, db) })
	router.DELETE("/products/:id", func(c *gin.Context) { handlers.DeleteProduct(c, db) })
	router.GET("/search-products/", func(c *gin.Context) { handlers.SearchAllProducts(c, db) })

	// Additional routes are defined here similarly...
}

// CORSMiddleware configures CORS (Cross-Origin Resource Sharing) headers to allow requests from specific origins.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
