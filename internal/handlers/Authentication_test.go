package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Set up your mock token service
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateTokenWithClaims(username, role string) (string, error) {
	args := m.Called(username, role)
	return args.String(0), args.Error(1)
}

// TestPostLogin tests the login method with the given username and password in different cases.
func TestPostLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		setupMock        func(db *gorm.DB, mock sqlmock.Sqlmock)
		username         string
		password         string
		expectedStatus   int
		expectedResponse string
		tokenError       error
	}{
		{
			name: "Successful login",
			setupMock: func(db *gorm.DB, mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"username", "password", "role"}).AddRow("user", hashedPassword, "user")
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WillReturnRows(rows)
			},
			username:         "user",
			password:         "correctpassword",
			expectedStatus:   http.StatusOK,
			expectedResponse: `"token":"mockedtoken"`,
			tokenError:       nil,
		},
		{
			name: "User not found",
			setupMock: func(db *gorm.DB, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WillReturnError(gorm.ErrRecordNotFound)
			},
			username:         "nonexistent",
			password:         "password",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `"error":"server error"`,
		},
		{
			name: "Incorrect password",
			setupMock: func(db *gorm.DB, mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"username", "password", "role"}).AddRow("user", hashedPassword, "user")
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WillReturnRows(rows)
			},
			username:         "user",
			password:         "wrongpassword",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: `"error":"authentication failed"`,
		},
		{
			name: "Token generation failure",
			setupMock: func(db *gorm.DB, mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"username", "password", "role"}).AddRow("user", hashedPassword, "user")
				mock.ExpectQuery("^SELECT \\* FROM \"users\" WHERE").WillReturnRows(rows)
			},
			username:         "user",
			password:         "correctpassword",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `"error":"could not generate token"`,
			tokenError:       errors.New("token generation failed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			mockTokenService := new(MockTokenService)
			if tc.tokenError == nil {
				mockTokenService.On("GenerateTokenWithClaims", tc.username, "user").Return("mockedtoken", nil)
			} else {
				mockTokenService.On("GenerateTokenWithClaims", tc.username, "user").Return("", tc.tokenError)
			}

			router.POST("/login", func(c *gin.Context) {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
				tc.setupMock(gormDB, mock)
				PostLogin(c, gormDB, mockTokenService)
			})

			bodyData := map[string]string{"username": tc.username, "password": tc.password}
			bodyBytes, _ := json.Marshal(bodyData)
			req, _ := http.NewRequest("POST", "/login", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.expectedResponse)
		})
	}
}
