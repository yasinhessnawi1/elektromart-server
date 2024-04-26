package tools

import (
	"strings"
)

/*
CheckString checks if the string is empty or not returns true if it is not
*/
func CheckString(stringToCheck string, maxLength int) bool {
	if stringToCheck == "" || len(stringToCheck) == 0 || len(stringToCheck) > maxLength {
		return false
	}
	return true
}

/*
CheckPassword checks if the password is less than 8 characters, contains a number, an uppercase letter, and a lowercase letter
*/
func CheckPassword(password string) bool {
	if len(password) < 8 || !strings.ContainsAny(password, "1234567890") ||
		!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") ||
		!strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") ||
		!strings.ContainsAny(password, "!@#$%^&*()_+|~-=\\`{}[]:\";'<>?,./") {
		return false
	}
	return true
}

/*
CheckInt checks if the int is less than 0, returns true if it is not
*/
func CheckInt(intToCheck int) bool {
	if intToCheck < 0 {
		return false
	}
	return true
}

/*
CheckFloat checks if the float is less than 0 returns true if it is not
*/
func CheckFloat(floatToCheck float64) bool {
	if floatToCheck < 0 {
		return false
	}
	return true
}

func CheckEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func CheckStatus(status string, i int) bool {
	status = strings.ToLower(status)
	validStatuses := []string{"pending", "shipped", "delivered", "returned", "cancelled", "refunded", "processing", "completed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func CheckDate(date string) bool {
	date = strings.TrimSpace(date)
	if len(date) != 10 || date[4] != '-' || date[7] != '-' {
		return false
	} else {
		for i, char := range date {
			if i == 4 || i == 7 {
				continue
			} else if char < '0' || char > '9' {
				return false
			}
		}
	}
	return true

}

func CheckPaymentMethod(method string) bool {
	method = strings.ToLower(method)
	validMethods := []string{"credit card", "debit card", "paypal", "cash", "check"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}
