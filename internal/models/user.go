package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	User_ID    uint32 `json:"user_id"`
	Username   string `gorm:"unique" json:"username"`
	Password   string `json:"password"`
	Email      string `gorm:"unique" json:"email"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Address    string `json:"address"`
}

type UserDB struct {
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
