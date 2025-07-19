package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rajvirsingh2/ascend-api/ai"
	"github.com/rajvirsingh2/ascend-api/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type QuestController struct {
	DB             *gorm.DB
	QuestGenerator ai.QuestGenerator
}

func NewQuestController(db *gorm.DB, questGen ai.QuestGenerator) *QuestController {
	return &QuestController{DB: db, QuestGenerator: questGen}
}

func (qc *QuestController) GenerateQuests(c *gin.Context) {
	var body struct {
		Goal string `json:"goal" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Goal is required"})
		return
	}

	user, _ := c.Get("user")
	currentUser := user.(models.User)

	var playerProfile models.PlayerProfile
	qc.DB.First(&playerProfile, "user_id = ?", currentUser.ID)

	generatedQuests, err := qc.QuestGenerator.GenerateQuests(c.Request.Context(), body.Goal, playerProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate quests"})
		return
	}

	var newQuests []models.Quest
	for _, q := range generatedQuests {
		newQuests = append(newQuests, models.Quest{
			PlayerProfileID: playerProfile.ID,
			Title:           q.Title,
			Description:     q.Description,
			Type:            models.QuestType(q.Type),
			XP:              q.XP,
			AttributeReward: q.AttributeReward,
		})
	}
	if err := qc.DB.Create(&newQuests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quests"})
		return
	}
	c.JSON(http.StatusOK, newQuests)
}

// GET (/quests)
func (qc *QuestController) GetActiveQuests(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)
	var playerProfile models.PlayerProfile
	qc.DB.First(&playerProfile, "user_id = ?", currentUser.ID)
	var activeQuests []models.Quest
	qc.DB.Where("player_profile_id = ? AND status=?", playerProfile.ID, models.StatusActive).Find(&activeQuests)
	c.JSON(http.StatusOK, activeQuests)
}

func (qc *QuestController) CompleteQuest(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)

	questId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest id"})
		return
	}

	var quest models.Quest
	var playerProfile models.PlayerProfile

	err = qc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&playerProfile, "user_id = ?", currentUser.ID).Error; err != nil {
			return err
		}
		if err := tx.Where("id=? AND player_profile_id = ? AND status=?", questId, playerProfile.ID, models.StatusActive).First(&quest).Error; err != nil {
			return fmt.Errorf("quest not found or not active")
		}

		quest.Status = models.StatusCompleted
		if err := tx.Save(&quest).Error; err != nil {
			return err
		}

		playerProfile.XP += quest.XP
		switch quest.AttributeReward {
		case "STRENGTH":
			playerProfile.Strength++
		case "AGILITY":
			playerProfile.Agility++
		case "INTELLIGENCE":
			playerProfile.Intelligence++
		case "VITALITY":
			playerProfile.Vitality++
		case "SENSE":
			playerProfile.Sense++
		}

		xpNextLevel := playerProfile.Level * 1000
		if playerProfile.XP >= xpNextLevel {
			playerProfile.Level++
			playerProfile.XP -= xpNextLevel
		}

		if err := tx.Save(&playerProfile).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, playerProfile)
}
