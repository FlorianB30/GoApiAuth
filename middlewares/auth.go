package middlewares

import (
	"auth-api-go/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware to Validate JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Debug: Check if middleware is executed
		fmt.Println("🔹 Middleware Executed")

		// Debug: Print received token
		fmt.Println("🔹 Raw Authorization Header:", tokenString)

		if tokenString == "" {
			fmt.Println("❌ Token missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("❌ Invalid token format:", tokenString)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			fmt.Println("❌ Token validation failed:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Debug: Print extracted claims
		fmt.Println("✅ Token claims:", claims)

		email, exists := claims["email"].(string)
		if !exists {
			fmt.Println("❌ Email not found in token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
			c.Abort()
			return
		}

		// Debug: Print extracted email
		fmt.Println("✅ Extracted Email:", email)

		// Set user email in context for further use
		c.Set("email", email)
		c.Next()
	}
}
