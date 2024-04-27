package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// OrderItem represents the order item model for an e-commerce transaction.
// It includes foreign keys to Order and Product, as well as Quantity and Subtotal to detail the item specifics.
type OrderItem struct {
	gorm.Model
	Order_ID   uint32  `json:"order_id"`
	Product_ID uint32  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Subtotal   float64 `json:"subtotal"`
}

// GetAllOrderItems retrieves all order items from the database.
// It returns a slice of OrderItem and an error if there is any issue during fetching.
func GetAllOrderItems(db *gorm.DB) ([]OrderItem, error) {
	var orderItems []OrderItem
	if err := db.Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

// SetOrderID validates and sets the Order_ID for an order item, ensuring the order exists.
// It returns true if the order exists and the ID is successfully set; otherwise, it returns false.
func (oi *OrderItem) SetOrderID(order_id uint32, db *gorm.DB) bool {
	if !OrderExists(db, order_id) {
		return false
	} else {
		oi.Order_ID = order_id
		return true
	}
}

// SetProductID validates and sets the Product_ID for an order item, ensuring the product exists.
// It returns true if the product exists and the ID is successfully set; otherwise, it returns false.
func (oi *OrderItem) SetProductID(product_id uint32, db *gorm.DB) bool {
	if !ProductExists(db, product_id) {
		return false
	} else {
		oi.Product_ID = product_id
		return true
	}
}

// SetQuantity validates and sets the quantity of an order item.
// It ensures the quantity is a positive integer before setting. Returns true if valid; otherwise false.
func (oi *OrderItem) SetQuantity(quantity int) bool {
	if !tools.CheckInt(quantity) {
		return false
	} else {
		oi.Quantity = quantity
		return true
	}
}

// SetSubtotal validates and sets the subtotal for an order item.
// It ensures the subtotal is a positive float before setting. Returns true if valid; otherwise false.
func (oi *OrderItem) SetSubtotal(subtotal float64) bool {
	if !tools.CheckFloat(subtotal) {
		return false
	} else {
		oi.Subtotal = subtotal
		return true
	}
}

// OrderItemExists checks if an order item exists in the database by its ID.
// It returns true if the order item is found, otherwise returns false.
func OrderItemExists(db *gorm.DB, id uint32) bool {
	var orderItem OrderItem
	if db.Where("id = ?", id).First(&orderItem).Error != nil {
		return false
	}
	return true
}

func SearchOrderItem(db *gorm.DB, searchParams map[string]interface{}) ([]OrderItem, error) {
	var orderItems []OrderItem
	query := db.Model(&OrderItem{})

	for key, value := range searchParams {
		switch key {
		case "order_id", "product_id":
			if numVal, ok := value.(int); ok {
				query = query.Where(key+" = ?", numVal)
			}
		case "quantity":
			// For numeric fields
			if numVal, ok := value.(int); ok {
				query = query.Where(key+" = ?", numVal)
			}
		case "subtotal":
			// For numeric fields
			if numVal, ok := value.(float64); ok {
				query = query.Where(key+" = ?", numVal)
			}
		}
	}

	if err := query.Find(&orderItems).Debug().Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}
