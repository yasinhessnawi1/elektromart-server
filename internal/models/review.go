package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Product_ID  uint32 `json:"product_id"`
	User_ID     uint32 `json:"user_id"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment"`
	Review_Date string `json:"review_date"`
}

func GetAllReviews(db *gorm.DB) ([]Review, error) {
	var reviews []Review
	if err := db.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *Review) SetProductID(product_id uint32, db *gorm.DB) bool {
	if !ProductExists(db, product_id) {
		return false
	} else {
		r.Product_ID = product_id
		return true
	}
}

func (r *Review) SetUserID(user_id uint32, db *gorm.DB) bool {
	if !UserExists(db, user_id) {
		return false
	} else {
		r.User_ID = user_id
		return true
	}
}

func (r *Review) SetRating(rating int) bool {
	if !tools.CheckRating(rating) {
		return false
	} else {
		r.Rating = rating
		return true
	}
}

func (r *Review) SetComment(comment string) bool {
	if !tools.CheckString(comment, 255) {
		return false
	} else {
		r.Comment = comment
		return true
	}
}

func (r *Review) SetReviewDate(review_date string) bool {
	if !tools.CheckDate(review_date) {
		return false
	} else {
		r.Review_Date = review_date
		return true
	}
}

func ReviewExists(db *gorm.DB, id uint32) bool {
	var review Review
	if db.Where("id = ?", id).First(&review).Error != nil {
		return false
	}
	return true
}

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
