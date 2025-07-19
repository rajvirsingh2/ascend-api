package models

import "gorm.io/gorm"

// Holds gamification data for a user
type PlayerProfile struct {
	gorm.Model
	UserId       uint `gorm:"unique; not null"`
	User         User `gorm:"foreignKey:UserId"`
	Level        int  `gorm:"default:1"`
	XP           int  `gorm:"default:0"`
	Strength     int  `gorm:"default:10"`
	Agility      int  `gorm:"default:10"`
	Intelligence int  `gorm:"default:10"`
	Vitality     int  `gorm:"default:10"`
	Sense        int  `gorm:"default:10"`
}
