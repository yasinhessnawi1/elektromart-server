package handlers

import (
	"E-Commerce_Website_Database/internal/models"
	"E-Commerce_Website_Database/internal/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

// GetShippingDetail retrieves a single shipping detail by its ID.
// It checks for the shipping detail's existence and validity of its data, then returns the shipping detail details or an error message.
// If the shipping detail is not found, it responds with an HTTP 404 Not Found status.
// If the shipping detail is found, it responds with an HTTP 200 OK status and the shipping detail details in JSON format.
func GetShippingDetail(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var shippingDetail models.ShippingDetails

	if err := db.Where("id = ?", id).First(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping Detail not found"})
		return
	}
	c.JSON(http.StatusOK, shippingDetail)
}

// GetShippingDetails retrieves all shipping details from the database.
// It returns a JSON response with a list of shipping details or an error message if the retrieval fails.
// If there are no shipping details in the database, it responds with an HTTP 404 Not Found status.
// If the retrieval is successful, it responds with an HTTP 200 OK status and the list of shipping details in JSON format.
func GetShippingDetails(c *gin.Context, db *gorm.DB) {
	shippingDetails, err := models.GetAllShippingDetails(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Shipping Details"})
		return
	}
	c.JSON(http.StatusOK, shippingDetails)
}

// SearchAllShippingDetails performs a search on shipping details based on provided query parameters.
// It constructs a search query dynamically and returns the matching shipping details or an appropriate error message.
// If no shipping details are found, it responds with an HTTP 404 Not Found status.
// If the search is successful, it responds with an HTTP 200 OK status and the list of shipping details in JSON format.
func SearchAllShippingDetails(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"order_id", "address", "shipping_date", "estimated_arrival", "status"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			switch field {
			case "order_id":
				if numVal, err := strconv.Atoi(cleanValue); err == nil {
					searchParams[field] = numVal
				}
			case "address", "shipping_date", "estimated_arrival":
				searchParams[field] = cleanValue
			case "status":
				searchParams[field] = strings.ToLower(cleanValue)
			default:
				searchParams[field] = cleanValue
			}
		}
	}

	shippingDetail, err := models.SearchShippingDetails(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Shipping Details", "details": err.Error()})
		return
	}

	if len(shippingDetail) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Shipping Details found"})
		return
	}

	c.JSON(http.StatusOK, shippingDetail)
}

// CreateShippingDetail creates a new shipping detail record in the database.
// It validates the incoming JSON data, creates a new shipping detail, and returns the newly created shipping detail or an error message.
// If the JSON data is invalid, it responds with an HTTP 400 Bad Request status.
// If the creation is successful, it responds with an HTTP 201 Created status and the created shipping detail in JSON format.
func CreateShippingDetail(c *gin.Context, db *gorm.DB) {
	var newShippingDetail models.ShippingDetails
	if err := c.ShouldBindJSON(&newShippingDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	shippingDetail := models.ShippingDetails{
		Order_ID:          newShippingDetail.Order_ID,
		Address:           newShippingDetail.Address,
		Shipping_Date:     newShippingDetail.Shipping_Date,
		Estimated_Arrival: newShippingDetail.Estimated_Arrival,
		Status:            newShippingDetail.Status,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkShippingDetail(shippingDetail, newShippingDetail, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Shipping Detail", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shippingDetail)
}

// UpdateShippingDetail updates an existing shipping detail record in the database.
// It validates the incoming JSON data, updates the shipping detail, and returns the updated shipping detail or an error message.
// If the JSON data is invalid, it responds with an HTTP 400 Bad Request status.
// If the update is successful, it responds with an HTTP 200 OK status and the updated shipping detail in JSON format.
func UpdateShippingDetail(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.ShippingDetailsExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping Detail not found"})
		return
	}

	var updatedShippingDetail models.ShippingDetails
	if err := c.ShouldBindJSON(&updatedShippingDetail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var shippingDetail models.ShippingDetails
	if err := db.Where("id = ?", id).First(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping Detail not found"})
		return
	}

	shippingDetail.Order_ID = updatedShippingDetail.Order_ID
	shippingDetail.Address = updatedShippingDetail.Address
	shippingDetail.Shipping_Date = updatedShippingDetail.Shipping_Date
	shippingDetail.Estimated_Arrival = updatedShippingDetail.Estimated_Arrival
	shippingDetail.Status = updatedShippingDetail.Status

	if failed, err := checkShippingDetail(shippingDetail, updatedShippingDetail, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipping detail", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shippingDetail)
}

// DeleteShippingDetail deletes a shipping detail record from the database.
// It checks for the existence of the shipping detail, deletes it, and responds with an appropriate status code.
// If the shipping detail does not exist, it responds with an HTTP 404 Not Found status.
// If the deletion is successful, it responds with an HTTP 204 No Content status.
func DeleteShippingDetail(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.ShippingDetailsExists(db, convertedId) {
		fmt.Println("Shipping Detail does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping Detail not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.ShippingDetails{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting Shipping Detail"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// checkShippingDetail validates the new shipping detail data against the existing shipping detail.
// It checks the order ID, address, shipping date, estimated arrival, and status fields for validity.
// If any field is invalid, it returns an error message and true, indicating a failed validation.
// If all fields are valid, it returns false and nil.
func checkShippingDetail(shippingDetail models.ShippingDetails, newShippingDetail models.ShippingDetails, db *gorm.DB) (bool, error) {
	switch true {
	case !shippingDetail.SetOrderID(newShippingDetail.Order_ID, db):
		return true, fmt.Errorf("invalid order id or not existing")
	case !shippingDetail.SetAddress(newShippingDetail.Address):
		return true, fmt.Errorf("invalid address")
	case !shippingDetail.SetShippingDate(newShippingDetail.Shipping_Date):
		return true, fmt.Errorf("shipping date is not expected")
	case !shippingDetail.SetEstimatedArrival(newShippingDetail.Estimated_Arrival):
		return true, fmt.Errorf("estimate date is not expected")
	case !shippingDetail.SetStatus(newShippingDetail.Status):
		return true, fmt.Errorf("status is not valid")
	}
	return false, nil
}
