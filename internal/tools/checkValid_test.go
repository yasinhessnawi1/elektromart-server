package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCheckString testing different cases of string inputs
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
