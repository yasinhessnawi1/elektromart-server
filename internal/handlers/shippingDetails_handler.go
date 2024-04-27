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

func GetShippingDetail(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var shippingDetail models.ShippingDetails

	if err := db.Where("id = ?", id).First(&shippingDetail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping Detail not found"})
		return
	}
	c.JSON(http.StatusOK, shippingDetail)
}

func GetShippingDetails(c *gin.Context, db *gorm.DB) {
	shippingDetail, err := models.GetAllShippingDetails(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Shipping Details"})
		return
	}
	c.JSON(http.StatusOK, shippingDetail)
}

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
