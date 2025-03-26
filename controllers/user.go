package controllers

import (
	"auth-api-go/config"
	"auth-api-go/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get User Profile (Protected Route)
func GetUserProfile(c *gin.Context) {
	// Extract user email from token
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing token"})
		return
	}

	// Find user in database
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Debugging print
	fmt.Println("🔹 Returning profile for:", user.Email)
	fmt.Println("test : ", user.Email, email)

	// Return user profile
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
