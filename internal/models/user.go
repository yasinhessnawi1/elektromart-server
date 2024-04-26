package models

import (
	"E-Commerce_Website_Database/internal/tools"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string `gorm:"unique" json:"username"`
	Password   string `json:"password"`
	Email      string `gorm:"unique" json:"email"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Address    string `json:"address"`
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) SetUsername(username string) bool {
	if !tools.CheckString(username, 255) {
		return false
	} else {
		u.Username = username
		return true
	}
}

func (u *User) SetPassword(password string) bool {
	if !tools.CheckPassword(password) {
		return false
	} else {
		u.Password = password
		return true
	}
}

func (u *User) SetEmail(email string) bool {
	if !tools.CheckEmail(email) {
		return false
	} else {
		u.Email = email
		return true
	}
}

func (u *User) SetFirstName(first_name string) bool {
	if !tools.CheckString(first_name, 255) {
		return false
	} else {
		u.First_Name = first_name
		return true
	}
}

func (u *User) SetLastName(last_name string) bool {
	if !tools.CheckString(last_name, 255) {
		return false
	} else {
		u.Last_Name = last_name
		return true
	}
}

func (u *User) SetAddress(address string) bool {
	if !tools.CheckString(address, 255) {
		return false
	} else {
		u.Address = address
		return true
	}
}

func UserExists(db *gorm.DB, id uint32) bool {
	var user User
	if db.First(&user, id).Error != nil {
		return false
	}
	return true
}

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
