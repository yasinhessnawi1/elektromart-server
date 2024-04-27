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

func GetReview(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var review models.Review

	if err := db.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}
	c.JSON(http.StatusOK, review)
}

func GetReviews(c *gin.Context, db *gorm.DB) {
	reviews, err := models.GetAllReviews(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving reviews"})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func SearchAllReviews(c *gin.Context, db *gorm.DB) {
	searchParams := map[string]interface{}{}

	for _, field := range []string{"product_id", "user_id", "rating", "comment", "review_date"} {
		if value := c.Query(field); value != "" {
			cleanValue := strings.TrimSpace(value)
			switch field {
			case "product_id", "user_id", "rating":
				if numVal, err := strconv.Atoi(cleanValue); err == nil {
					searchParams[field] = numVal
				}
			case "comment", "review_date":
				searchParams[field] = cleanValue
			default:
				searchParams[field] = cleanValue
			}
		}
	}

	reviews, err := models.SearchReview(db, searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reviews", "details": err.Error()})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No reviews found"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func CreateReview(c *gin.Context, db *gorm.DB) {
	var newReview models.Review
	if err := c.ShouldBindJSON(&newReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	review := models.Review{
		Product_ID:  newReview.Product_ID,
		User_ID:     newReview.User_ID,
		Rating:      newReview.Rating,
		Comment:     newReview.Comment,
		Review_Date: newReview.Review_Date,
		Model: gorm.Model{
			ID: uint(tools.GenerateUUID()),
		},
	}

	if failed, err := checkReview(review, newReview, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func UpdateReview(c *gin.Context, db *gorm.DB) {
	id := tools.ConvertStringToUint(c.Param("id"))

	if !models.ReviewExists(db, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	var updatedReview models.Review
	if err := c.ShouldBindJSON(&updatedReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	var review models.Review
	if err := db.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping Detail not found"})
		return
	}

	review.Product_ID = updatedReview.Product_ID
	review.User_ID = updatedReview.User_ID
	review.Rating = updatedReview.Rating
	review.Comment = updatedReview.Comment
	review.Review_Date = updatedReview.Review_Date

	if failed, err := checkReview(review, updatedReview, db); failed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
		return
	}

	if err := db.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

func DeleteReview(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	convertedId := tools.ConvertStringToUint(id)

	if !models.ReviewExists(db, convertedId) {
		fmt.Println("Review does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if err := db.Unscoped().Where("id = ?", convertedId).Delete(&models.Review{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting Review"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func checkReview(review models.Review, newReview models.Review, db *gorm.DB) (bool, error) {
	switch true {
	case !review.SetProductID(newReview.Product_ID, db):
		return true, fmt.Errorf("invalid Product id or not existing")
	case !review.SetUserID(newReview.User_ID, db):
		return true, fmt.Errorf("invalid user id or not existing")
	case !review.SetRating(newReview.Rating):
		return true, fmt.Errorf("the rate must be between 1 and 5(best)")
	case !review.SetComment(newReview.Comment):
		return true, fmt.Errorf("the comment is not valid")
	case !review.SetReviewDate(newReview.Review_Date):
		return true, fmt.Errorf("review date is not expected")
	}
	return false, nil
}
