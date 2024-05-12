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

// GetReview retrieves a single review by its ID.
// It checks for the review's existence and validity of its data, then returns the review details or an error message.
// If the review is not found, it responds with an HTTP 404 Not Found status.
// If the review is found, it responds with an HTTP 200 OK status and the review details in JSON format.
func GetReview(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var review models.Review

	if err := db.Where("id = ?", id).First(&review).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}
	c.JSON(http.StatusOK, review)
}

// GetReviews retrieves all reviews from the database.
// It returns a JSON response with a list of reviews or an error message if the retrieval fails.
// If there are no reviews in the database, it responds with an HTTP 404 Not Found status.
// If the retrieval is successful, it responds with an HTTP 200 OK status and the list of reviews in JSON format.
func GetReviews(c *gin.Context, db *gorm.DB) {
	reviews, err := models.GetAllReviews(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving reviews"})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

// SearchAllReviews performs a search on reviews based on provided query parameters.
// It constructs a search query dynamically and returns the matching reviews or an appropriate error message.
// If no reviews are found, it responds with an HTTP 404 Not Found status.
// If the search is successful, it responds with an HTTP 200 OK status and the list of reviews in JSON format.
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

// CreateReview adds a new review to the database.
// It validates the review data and responds with the newly created review or an error message.
// If the review data is invalid, it responds with an HTTP 400 Bad Request status.
// If the creation is successful, it responds with an HTTP 201 Created status and the created review in JSON format.
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

// UpdateReview updates an existing review in the database.
// It validates the updated review data and responds with the updated review or an error message.
// If the review data is invalid, it responds with an HTTP 400 Bad Request status.
// If the update is successful, it responds with an HTTP 200 OK status and the updated review in JSON format.
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

// DeleteReview removes a review from the database.
// It checks for the review's existence and responds with an appropriate status code.
// If the review is not found, it responds with an HTTP 404 Not Found status.
// If the deletion is successful, it responds with an HTTP 204 No Content status.
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

// checkReview validates the review data before creating or updating a review.
// It checks the product ID, user ID, rating, comment, and review date for validity.
// If any of the data is invalid, it returns an error message and true, indicating a failure.
// If all data is valid, it returns false and nil, indicating success.
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
