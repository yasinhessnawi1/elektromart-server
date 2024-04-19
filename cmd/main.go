package main

import (
	"E-Commerce_Website_Database/internal/config"
	"E-Commerce_Website_Database/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()

	// Correct MySQL connection string format
	dsn := "root:@tcp(localhost:8080)/eCommerce?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := gin.Default()

	setupRoutes(router, db)

	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
func setupRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to ElectroMart API"})
	})
	router.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
	router.GET("/orders", func(c *gin.Context) { handlers.GetOrders(c, db) })
	router.GET("/order_items", func(c *gin.Context) { handlers.GetOrderItems(c, db) })
	router.GET("/payments", func(c *gin.Context) { handlers.GetPayments(c, db) })
}
