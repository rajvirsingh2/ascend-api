package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rajvirsingh2/ascend-api/config"
	"github.com/rajvirsingh2/ascend-api/models"
	"gorm.io/gorm"
	"net/http"
)

func GetProfile(c *gin.Context) {
	fmt.Println("📩 /api/v1/profile hit")

	user, exists := c.Get("user")
	fmt.Println("🔍 context user exists?", exists)

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Player"})
		return
	}

	currentUser := user.(models.User)
	fmt.Println("👤 currentUser:", currentUser.ID, currentUser.Email)

	var playerProfile models.PlayerProfile
	result := config.DB.Preload("User").Where("user_id=?", currentUser.ID).Take(&playerProfile)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("⚠️ Player profile not found for user", currentUser.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Player profile not found"})
		return
	} else if result.Error != nil {
		fmt.Println("❌ DB error:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	fmt.Println("✅ Player profile:", playerProfile.ID)
	c.JSON(http.StatusOK, playerProfile)
}
