package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// TestGetAllReviews tests the retrieval of all reviews from the database.
func TestGetAllReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "rating", "comment", "review_date"}).
		AddRow(1, 1, 1, 5, "Great product", "2021-04-21").
		AddRow(2, 1, 2, 4, "Good, but expensive", "2021-04-22")
	mock.ExpectQuery("^SELECT \\* FROM \"reviews\"").WillReturnRows(rows)

	reviews, err := GetAllReviews(gormDB)
	assert.NoError(t, err)
	assert.Len(t, reviews, 2, "Should fetch two reviews")
}

// TestReview_SetProductID tests setting the product ID after verifying the product exists.
func TestReview_SetProductID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	review := Review{}
	result := review.SetProductID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"products\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = review.SetProductID(99, gormDB)
	assert.False(t, result)
}

// TestReview_SetUserID tests setting the user ID after verifying the user exists.
func TestReview_SetUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	review := Review{}
	result := review.SetUserID(1, gormDB)
	assert.True(t, result)

	mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WithArgs(99, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	result = review.SetUserID(99, gormDB)
	assert.False(t, result)
}

// TestReview_SetRating tests setting the rating after validating it with predefined rules.
func TestReview_SetRating(t *testing.T) {
	review := Review{}
	assert.False(t, review.SetRating(-1), "Rating should be invalid because it is negative")
	assert.True(t, review.SetRating(5), "Rating should be valid")
}

// TestReview_SetComment tests setting the comment after validating its length.
func TestReview_SetComment(t *testing.T) {
	review := Review{}
	assert.False(t, review.SetComment(string(make([]byte, 256))), "Comment should be invalid due to length")
	assert.True(t, review.SetComment("Great product!"), "Comment should be valid")
}

// TestReview_SetReviewDate tests setting the review date after validating it as a valid date string.
func TestReview_SetReviewDate(t *testing.T) {
	review := Review{}
	assert.False(t, review.SetReviewDate("20210421"), "Review date should be invalid due to format")
	assert.True(t, review.SetReviewDate("2021-04-21"), "Review date should be valid")
}

// TestReviewExists tests checking if a specific review exists in the database by its ID.
func TestReviewExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	mock.ExpectQuery("^SELECT \\* FROM \"reviews\" WHERE").WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	assert.True(t, ReviewExists(gormDB, 1), "Review should exist")

	mock.ExpectQuery("^SELECT \\* FROM \"reviews\" WHERE").WithArgs(2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	assert.False(t, ReviewExists(gormDB, 2), "Review should not exist")
}

// TestSearchReview tests the search functionality based on provided parameters.
func TestSearchReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "product_id", "user_id", "rating", "comment", "review_date"}).
		AddRow(1, 1, 1, 5, "Great product", "2021-04-21")
	mock.ExpectQuery("^SELECT \\* FROM \"reviews\" WHERE").
		WithArgs(1, 5).
		WillReturnRows(rows)

	searchParams := map[string]interface{}{
		"product_id": 1,
		"rating":     5,
	}
	reviews, err := SearchReview(gormDB, searchParams)
	assert.NoError(t, err)
	assert.Len(t, reviews, 1, "Should find one matching review")
}
