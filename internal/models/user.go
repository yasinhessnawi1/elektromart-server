package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

// User represents the user entity in the database.
// It includes essential fields like Username, Password, Email, along with personal details such as First Name, Last Name, and Address.
type User struct {
	gorm.Model
	Username   string `gorm:"unique" json:"username"`
	Password   string `json:"password"`
	Email      string `gorm:"unique" json:"email"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Address    string `json:"address"`
	Mobile     string `json:"mobile"`
	Role       string `json:"role"`
}

// GetAllUsers retrieves all users from the database.
// Returns a slice of User or an error if the fetch fails.
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) SetRole(role string) bool {
	if !tools.CheckRole(role) {
		return false
	} else {
		u.Role = role
		return true
	}

}

// SetUsername sets the username for the user after validating its uniqueness and length.
// Returns true if the username is within the allowed length and unique, otherwise false.
func (u *User) SetUsername(username string) bool {
	if !tools.CheckString(username, 255) {
		return false
	} else {
		u.Username = username
		return true
	}
}

// SetUsername sets the username for the user after validating its uniqueness and length.
// Returns true if the username is within the allowed length and unique, otherwise false.
func (u *User) SetPhone(phone string) bool {
	if !tools.CheckPhone(phone, 11) {
		return false
	} else {
		u.Mobile = phone
		return true
	}
}

// SetPassword sets the password for the user after ensuring it meets security standards.
// Returns true if the password is considered secure, otherwise false.
func (u *User) SetPassword(password string) bool {
	if !tools.CheckPassword(password) {
		return false
	} else {
		u.Password = password
		return true
	}
}

// SetEmail sets the email for the user after validating its format and uniqueness.
// Returns true if the email is valid and unique, otherwise false.
func (u *User) SetEmail(email string) bool {
	if !tools.CheckEmail(email) {
		return false
	} else {
		u.Email = email
		return true
	}
}

// SetFirstName sets the first name of the user after validating its length.
// Returns true if the first name is within the allowed length, otherwise false.
func (u *User) SetFirstName(first_name string) bool {
	if !tools.CheckString(first_name, 255) {
		return false
	} else {
		u.First_Name = first_name
		return true
	}
}

// SetLastName sets the last name of the user after validating its length.
// Returns true if the last name is within the allowed length, otherwise false.
func (u *User) SetLastName(last_name string) bool {
	if !tools.CheckString(last_name, 255) {
		return false
	} else {
		u.Last_Name = last_name
		return true
	}
}

// SetAddress sets the address for the user after validating its length.
// Returns true if the address is within the allowed length, otherwise false.
func (u *User) SetAddress(address string) bool {
	if !tools.CheckString(address, 255) {
		return false
	} else {
		u.Address = address
		return true
	}
}

// UserExists checks if a specific user exists in the database by their ID.
// Returns true if the user exists, otherwise false.
func UserExists(db *gorm.DB, id uint32) bool {
	var user User
	if db.First(&user, id).Error != nil {
		return false
	}
	return true
}

// SearchUsers performs a search based on given search parameters.
// It filters users based on criteria like username, email, first name, last name, and address.
// Returns a slice of users that match the criteria or an error if the search fails.
func SearchUsers(db *gorm.DB, searchParams map[string]interface{}) ([]User, error) {
	var users []User
	query := db.Model(&User{})

	for key, value := range searchParams {
		if strVal, ok := value.(string); ok {
			query = query.Where(key+" LIKE ?", "%"+strVal+"%")
		}
	}

	if err := query.Find(&users).Debug().Error; err != nil {
		return nil, err
	}
	return users, nil
}
