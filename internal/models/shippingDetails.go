package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
	"strings"
)

// ShippingDetails represents the shipping details model for an e-commerce transaction.
// It includes fields for Order_ID, Address, Shipping_Date, Estimated_Arrival, and Status.
type ShippingDetails struct {
	gorm.Model
	Order_ID          uint32 `json:"order_id"`
	Address           string `json:"address"`
	Shipping_Date     string `json:"shipping_date"`
	Estimated_Arrival string `json:"estimated_arrival"`
	Status            string `json:"status"`
}

// GetAllShippingDetails retrieves all shipping details from the database.
// It returns a slice of ShippingDetails or an error if the fetch fails.
func GetAllShippingDetails(db *gorm.DB) ([]ShippingDetails, error) {
	var shippingDetails []ShippingDetails
	if err := db.Find(&shippingDetails).Error; err != nil {
		return nil, err
	}
	return shippingDetails, nil
}

// SetOrderID sets the order ID for the shipping details after verifying the existence of the order.
// Returns true if the order exists and the ID is set; otherwise, it returns false.
func (s *ShippingDetails) SetOrderID(order_id uint32, db *gorm.DB) bool {
	if !OrderExists(db, order_id) {
		return false
	} else {
		s.Order_ID = order_id
		return true
	}
}

// SetAddress sets the address for the shipping details after validating its length.
// Returns true if the address is within the allowed length, otherwise false.
func (s *ShippingDetails) SetAddress(address string) bool {
	if !tools.CheckString(address, 255) {
		return false
	} else {
		s.Address = address
		return true
	}
}

// SetShippingDate sets the shipping date for the shipping details after validating the date format.
// Returns true if the date is valid, otherwise false.
func (s *ShippingDetails) SetShippingDate(shipping_date string) bool {
	if !tools.CheckDate(shipping_date) {
		return false
	} else {
		s.Shipping_Date = shipping_date
		return true
	}
}

// SetEstimatedArrival sets the estimated arrival date for the shipping details after validating the date format.
// Returns true if the date is valid, otherwise false.
func (s *ShippingDetails) SetEstimatedArrival(estimated_arrival string) bool {
	if !tools.CheckDate(estimated_arrival) {
		return false
	} else {
		s.Estimated_Arrival = estimated_arrival
		return true
	}
}

// SetStatus sets the status for the shipping details after validating its length.
// Returns true if the status is within the allowed length, otherwise false.
func (s *ShippingDetails) SetStatus(status string) bool {
	if !tools.CheckStatus(status, 255) {
		return false
	} else {
		s.Status = status
		return true
	}
}

// ShippingDetailsExists checks if a shipping details record exists in the database by its ID.
// It returns true if the shipping details record is found, otherwise returns false.
func ShippingDetailsExists(db *gorm.DB, id uint32) bool {
	var shippingDetails ShippingDetails
	if db.Where("id = ?", id).First(&shippingDetails).Error != nil {
		return false
	}
	return true
}

// SearchShippingDetails adds a new shipping details record to the database.
// It returns the new shipping details record or an error if the operation fails.
// The shipping details record is created using the provided ShippingDetails struct.
func SearchShippingDetails(db *gorm.DB, searchParams map[string]interface{}) ([]ShippingDetails, error) {
	var shippingDetails []ShippingDetails
	query := db.Model(&ShippingDetails{})

	for key, value := range searchParams {
		valueStr, isString := value.(string)
		switch key {
		case "order_id":
			if numVal, ok := value.(int); ok {
				query = query.Where(key+" = ?", numVal)
			}
		case "address":
			if isString {
				query = query.Where(key+" LIKE ?", "%"+valueStr+"%")
			}
		case "shipping_date", "estimated_arrival":
			if isString {
				query = query.Where(key+" = ?", valueStr)
			}
		case "status":
			if isString {
				query = query.Where(key+" LIKE ?", "%"+strings.ToLower(valueStr)+"%")
			}
		}
	}

	if err := query.Find(&shippingDetails).Debug().Error; err != nil {
		return nil, err
	}
	return shippingDetails, nil
}
