package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCheckString testing different cases of string inputs
// It checks if the string is empty, valid, or exceeds the maximum length
// Returns true if the string is within the valid range, otherwise false.
func TestCheckString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected bool
	}{
		{"Empty string", "", 10, false},
		{"Valid string", "Hello", 10, true},
		{"Invalid Max Length", "Hello Hello!", 10, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckString(test.input, test.maxLen)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckPassword testing different cases for password validation
// It checks if the password is at least 8 characters long and includes at least one number, one uppercase letter, one lowercase letter, and one special character
// Returns true if the password meets these criteria, otherwise false.
func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Too Short", "Pass", false},
		{"No Number", "Password!", false},
		{"No Upper", "password1!", false},
		{"Valid Pass", "Password1!", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckPassword(test.password)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckInt testing different cases for integer numbers validation
// Returns true if the integer is 0 or positive, otherwise false.
func TestCheckInt(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		expected bool
	}{
		{"Negative", -1, false},
		{"Valid number", 2, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckInt(test.number)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckRating testing different cases for rating validation
// Returns true if the integer is between 0 and 5, otherwise false.
func TestCheckRating(t *testing.T) {
	tests := []struct {
		name     string
		rating   int
		expected bool
	}{
		{"Negative", -1, false},
		{"Too High number", 6, false},
		{"Valid Lower number", 0, true},
		{"Valid Upper number", 5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckRating(test.rating)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckFloat testing different cases for floating numbers validation
// Returns true if the number is 0.0 or greater, otherwise false.
func TestCheckFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"Negative", -1.0, false},
		{"Zero", 0.0, true},
		{"Positive", 1.0, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckFloat(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckEmail testing different cases for email validation
// Returns true if the string looks like an email address, otherwise false.
func TestCheckEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid Email", "example@example.com", true},
		{"No @ Symbol", "example.com", false},
		{"No Dot", "example@example", false},
		{"Empty String", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckEmail(test.email)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckStatus testing different cases for status validation
// Returns true if the status is valid, otherwise false.
func TestCheckStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Valid Status", "pending", true},
		{"Invalid Status", "waiting", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckStatus(test.status, 20)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckDate testing different cases for date validation
// Returns true if the format is correct, otherwise false.
func TestCheckDate(t *testing.T) {
	tests := []struct {
		name     string
		date     string
		expected bool
	}{
		{"Valid Date", "2023-04-01", true},
		{"Invalid Date", "2023/04/01", false},
		{"Invalid Date", "20ab-cd-0e", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckDate(test.date)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckPaymentMethod testing different cases for payment validation
// Returns true if the method is valid, otherwise false.
func TestCheckPaymentMethod(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected bool
	}{
		{"Valid Method", "credit card", true},
		{"Invalid Method", "crypto", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckPaymentMethod(test.method)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestCheckRole testing different cases for role validation
// Returns true if the role is valid, otherwise false.
func TestCheckRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"Valid Role", "admin", true},
		{"Invalid Role", "superuser", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CheckRole(test.role)
			assert.Equal(t, test.expected, result)
		})
	}
}
