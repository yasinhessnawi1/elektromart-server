package main

import (
	"E-Commerce_Website_Database/internal/config"
	"E-Commerce_Website_Database/internal/handlers"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

// main initializes the application, sets up database connections, and starts the HTTP server.
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
	// Here you should use Query Param Like :search-users/?username={The username}  or search-users/?email={The email}
	//`or by first name , last name , or address`.
	router.GET("/search-users/", func(c *gin.Context) { handlers.SearchAllUsers(c, db) })

	router.GET("/shippingDetails", func(c *gin.Context) { handlers.GetShippingDetails(c, db) })
	router.GET("/shippingDetails/:id", func(c *gin.Context) { handlers.GetShippingDetail(c, db) })
	router.POST("/shippingDetails", func(c *gin.Context) { handlers.CreateShippingDetail(c, db) })
	router.PUT("/shippingDetails/:id", func(c *gin.Context) { handlers.UpdateShippingDetail(c, db) })
	router.DELETE("/shippingDetails/:id", func(c *gin.Context) { handlers.DeleteShippingDetail(c, db) })
	// Here you should use Query Param Like :search-shippingDetails/?order_id={exist ID}  or search-shippingDetails/?address={The address}
	//`or by status`.
	router.GET("/search-shippingDetails/", func(c *gin.Context) { handlers.SearchAllShippingDetails(c, db) })

	router.GET("/reviews", func(c *gin.Context) { handlers.GetReviews(c, db) })
	router.GET("/reviews/:id", func(c *gin.Context) { handlers.GetReview(c, db) })
	router.POST("/reviews", func(c *gin.Context) { handlers.CreateReview(c, db) })
	router.PUT("/reviews/:id", func(c *gin.Context) { handlers.UpdateReview(c, db) })
	router.DELETE("/reviews/:id", func(c *gin.Context) { handlers.DeleteReview(c, db) })
	// Here you should use Query Param Like :search-reviews/?product_id={exist ID}  or search-reviews/?comment={The comment}
	//`or by rating, user_id , review_date`.
	router.GET("/search-reviews/", func(c *gin.Context) { handlers.SearchAllReviews(c, db) })

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
	// Here you should use Query Param Like :search-brands/?name={The name}  or search-brands/?description={The description}
	router.GET("/search-brands/", func(c *gin.Context) { handlers.SearchAllBrands(c, db) })

	router.GET("/categories", func(c *gin.Context) { handlers.GetCategories(c, db) })
	router.GET("/categories/:id", func(c *gin.Context) { handlers.GetCategory(c, db) })
	router.POST("/categories", func(c *gin.Context) { handlers.CreateCategory(c, db) })
	router.PUT("/categories/:id", func(c *gin.Context) { handlers.UpdateCategory(c, db) })
	router.DELETE("/categories/:id", func(c *gin.Context) { handlers.DeleteCategory(c, db) })
	// Here you should use Query Param Like :search-categories/?name={The name}  or search-categories/?description={The description}
	router.GET("/search-categories/", func(c *gin.Context) { handlers.SearchAllCategories(c, db) })

	router.GET("/orders", func(c *gin.Context) { handlers.GetOrders(c, db) })
	router.GET("/orders/:id", func(c *gin.Context) { handlers.GetOrder(c, db) })
	router.POST("/orders", func(c *gin.Context) { handlers.CreateOrder(c, db) })
	router.PUT("/orders/:id", func(c *gin.Context) { handlers.UpdateOrder(c, db) })
	router.DELETE("/orders/:id", func(c *gin.Context) { handlers.DeleteOrder(c, db) })
	// Here you should use Query Param Like :search-orders/?user_id={exist ID}  or search-orders/?total_amount={The amount}
	//`or by status`.
	router.GET("/search-orders/", func(c *gin.Context) { handlers.SearchAllOrders(c, db) })

	router.GET("/orderItems", func(c *gin.Context) { handlers.GetOrderItems(c, db) })
	router.GET("/orderItems/:id", func(c *gin.Context) { handlers.GetOrderItem(c, db) })
	router.POST("/orderItems", func(c *gin.Context) { handlers.CreateOrderItem(c, db) })
	router.PUT("/orderItems/:id", func(c *gin.Context) { handlers.UpdateOrderItem(c, db) })
	router.DELETE("/orderItems/:id", func(c *gin.Context) { handlers.DeleteOrderItem(c, db) })
	// Here you should use Query Param Like :search-orderItems/?order_id={the order id}  or search-orderItems/?quantity={The quantity}
	//`or by product id `.
	router.GET("/search-orderItems/", func(c *gin.Context) { handlers.SearchAllOrderItems(c, db) })

	router.GET("/payments", func(c *gin.Context) { handlers.GetPayments(c, db) })
	router.GET("/payments/:id", func(c *gin.Context) { handlers.GetPayment(c, db) })
	router.POST("/payments", func(c *gin.Context) { handlers.CreatePayment(c, db) })
	router.PUT("/payments/:id", func(c *gin.Context) { handlers.UpdatePayment(c, db) })
	router.DELETE("/payments/:id", func(c *gin.Context) { handlers.DeletePayment(c, db) })
	// Here you should use Query Param Like :search-payments/?payment_method={cash}  or search-payments/?amount={The amount}
	//`or by order id `.
	router.GET("/search-payments/", func(c *gin.Context) { handlers.SearchAllPayments(c, db) })
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Request received")
		c.Next()
		fmt.Println("Headers:", c.Writer.Header())
	}
}
