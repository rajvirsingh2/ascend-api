package models

import (
	"gorm.io/gorm"
	"time"
)

type QuestStatus string

const (
	StatusActive    QuestStatus = "ACTIVE"
	StatusCompleted QuestStatus = "COMPLETED"
)

type QuestType string

const (
	TypeDaily  QuestType = "DAILY"
	TypeWeekly QuestType = "WEEKLY"
)

// A single task for the player to complete.
type Quest struct {
	gorm.Model
	PlayerProfileID uint   `gorm:"not null"`
	Title           string `gorm:"type:varchar(255)"`
	Description     string
	Type            QuestType   `gorm:"type:varchar(10)"`
	Status          QuestStatus `gorm:"type:varchar(10);default:'ACTIVE'"`
	XP              int
	AttributeReward string `gorm:"type:varchar(20)"`
	CompletedAt     *time.Time
}
