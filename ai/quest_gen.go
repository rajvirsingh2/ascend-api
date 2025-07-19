package ai

import (
	"context"
	"github.com/rajvirsingh2/ascend-api/models"
)

type GeneratedQuest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Type            string `json:"type"`
	XP              int    `json:"xp"`
	AttributeReward string `json:"attribute_reward"`
}

type QuestGenerator interface {
	GenerateQuests(ctx context.Context, goal string, profile models.PlayerProfile) ([]GeneratedQuest, error)
}
