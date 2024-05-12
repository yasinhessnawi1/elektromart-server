package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// Review represents the review model for products.
// It includes fields for Product_ID, User_ID, Rating, Comment, and Review_Date.
type Review struct {
	gorm.Model
	Product_ID  uint32 `json:"product_id"`
	User_ID     uint32 `json:"user_id"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment"`
	Review_Date string `json:"review_date"`
}

// GetAllReviews retrieves all reviews from the database.
// If the fetch is successful, it returns the list of reviews, otherwise an error message.
func GetAllReviews(db *gorm.DB) ([]Review, error) {
	var reviews []Review
	if err := db.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// SetProductID sets the product ID for the review after verifying the existence of the product.
// Returns true if the product exists and the ID is set; otherwise, it returns false.
func (r *Review) SetProductID(product_id uint32, db *gorm.DB) bool {
	if !ProductExists(db, product_id) {
		return false
	} else {
		r.Product_ID = product_id
		return true
	}
}

// SetUserID sets the user ID for the review after verifying the existence of the user.
// Returns true if the user exists and the ID is set; otherwise, it returns false.
func (r *Review) SetUserID(user_id uint32, db *gorm.DB) bool {
	if !UserExists(db, user_id) {
		return false
	} else {
		r.User_ID = user_id
		return true
	}
}

// SetRating sets the rating for the review.
// It validates the rating to ensure it meets certain criteria and returns true if valid.
func (r *Review) SetRating(rating int) bool {
	if !tools.CheckRating(rating) {
		return false
	} else {
		r.Rating = rating
		return true
	}
}

// SetComment sets the comment for the review.
// It validates the comment to ensure it meets certain criteria and returns true if valid.
func (r *Review) SetComment(comment string) bool {
	if !tools.CheckString(comment, 255) {
		return false
	} else {
		r.Comment = comment
		return true
	}
}

// SetReviewDate sets the review date for the review.
// It validates the date format and returns true if the date is valid; otherwise, it returns false.
func (r *Review) SetReviewDate(review_date string) bool {
	if !tools.CheckDate(review_date) {
		return false
	} else {
		r.Review_Date = review_date
		return true
	}
}

// ReviewExists checks if a review exists in the database by its ID.
// It returns true if the review is found, otherwise returns false.
func ReviewExists(db *gorm.DB, id uint32) bool {
	var review Review
	if db.Where("id = ?", id).First(&review).Error != nil {
		return false
	}
	return true
}

// SearchReview retrieves reviews from the database based on the search parameters provided.
// It returns a slice of reviews that match the criteria or an error if the search fails.
func SearchReview(db *gorm.DB, searchParams map[string]interface{}) ([]Review, error) {
	var reviews []Review
	query := db.Model(&Review{})

	for key, value := range searchParams {
		valueStr, isString := value.(string)
		switch key {
		case "product_id", "user_id":
			if numVal, ok := value.(int); ok {
				query = query.Where(key+" = ?", numVal)
			}
		case "rating":
			// For numeric fields
			if numVal, ok := value.(int); ok {
				query = query.Where(key+" = ?", numVal)
			}
		case "comment":
			if isString {
				query = query.Where(key+" LIKE ?", "%"+valueStr+"%")
			}
		case "review_date":
			if isString {
				query = query.Where(key+" = ?", valueStr)
			}
		}
	}

	if err := query.Find(&reviews).Debug().Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
