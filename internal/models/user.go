package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	BaseModel
	Username  string `gorm:"unique"`
	Password  string
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Address   string
}
type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
