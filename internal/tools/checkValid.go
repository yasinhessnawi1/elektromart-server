package tools

import (
	"strings"
)

// CheckString validates the length of a string, ensuring it is not empty and does not exceed the specified maxLength.
// Returns true if the string is within the valid range, otherwise false.
func CheckString(stringToCheck string, maxLength int) bool {
	return stringToCheck != "" && len(stringToCheck) <= maxLength
}

// CheckPassword ensures that a password meets specific security criteria:
// It must be at least 8 characters long and include at least one number, one uppercase letter, one lowercase letter, and one special character.
// Returns true if the password meets these criteria, otherwise false.
func CheckPassword(password string) bool {
	hasMinLength := len(password) >= 8
	hasDigit := strings.ContainsAny(password, "1234567890")
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+|~-=\\`{}[]:\";'<>?,./")

	return hasMinLength && hasDigit && hasUpper && hasLower && hasSpecial
}

// CheckInt verifies if an integer is non-negative.
// Returns true if the integer is 0 or positive, otherwise false.
func CheckInt(intToCheck int) bool {
	return intToCheck >= 0
}

// CheckFloat checks if a floating-point number is non-negative.
// Returns true if the number is 0.0 or greater, otherwise false.
func CheckFloat(floatToCheck float64) bool {
	return floatToCheck >= 0.0
}

// CheckEmail verifies if a string contains basic elements that could constitute a valid email address:
// It must contain an '@' character and at least one dot '.'.
// Returns true if the string looks like an email address, otherwise false.
func CheckEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// CheckStatus validates if a given status string is one of the predefined valid statuses.
// Returns true if the status is valid, otherwise false.
func CheckStatus(status string, maxLength int) bool {
	validStatuses := []string{"pending", "shipped", "delivered", "returned", "cancelled", "refunded", "processing", "completed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// CheckDate validates a date string format to be YYYY-MM-DD.
// Returns true if the format is correct, otherwise false.
func CheckDate(date string) bool {
	date = strings.TrimSpace(date)
	if len(date) != 10 || date[4] != '-' || date[7] != '-' {
		return false
	}
	for i, char := range date {
		if i == 4 || i == 7 {
			continue
		} else if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// CheckPaymentMethod validates if a payment method is one of the predefined valid methods.
// Returns true if the method is valid, otherwise false.
func CheckPaymentMethod(method string) bool {
	validMethods := []string{"credit card", "debit card", "paypal", "cash", "check"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}
