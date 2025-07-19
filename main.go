package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rajvirsingh2/ascend-api/ai"
	"github.com/rajvirsingh2/ascend-api/config"
	"github.com/rajvirsingh2/ascend-api/controller"
	"github.com/rajvirsingh2/ascend-api/middleware"
	"github.com/rajvirsingh2/ascend-api/models"
	"log"
)

func init() {
	config.ConnectDB()
}

func main() {
	err := config.DB.AutoMigrate(&models.User{}, &models.PlayerProfile{}, &models.Quest{})
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	fmt.Println("✅ Database migrated")

	questGenerator, err := ai.NewGeminiAdapter()
	if err != nil {
		log.Fatalf("failed to create gemini adapter: %v", err)
	}
	fmt.Println("✅ AI Quest Generator initialized")

	questController := controller.NewQuestController(config.DB, questGenerator)
	r := gin.Default()

	//Public Routes
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", controller.Register)
		authRoute.POST("/login", controller.Login)
	}

	//Private Route
	api := r.Group("/api/v1")
	api.Use(middleware.RequireAuth)
	{
		api.GET("/profile", controller.GetProfile)
		questRoutes := api.Group("/quests")
		{
			questRoutes.POST("/generate", questController.GenerateQuests)
			questRoutes.GET("", questController.GetActiveQuests)
			questRoutes.POST("/:id/complete", questController.CompleteQuest)
		}
	}
	r.Run("0.0.0.0:8000")
}
