package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rajvirsingh2/ascend-api/config"
	"github.com/rajvirsingh2/ascend-api/models"
	"net/http"
	"os"
	"strings"
	"time"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		authReader := c.GetHeader("Authorization")
		if authReader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authReader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString = parts[1]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID := uint(userIDFloat)

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Println("✅ User fetched from DB:", user.ID, user.Email)
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
