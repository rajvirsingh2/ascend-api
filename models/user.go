package models

import "gorm.io/gorm"

// User Authenticated for application
type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"notnull"`
}
