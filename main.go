package main

import (
	"auth-api-go/config"
	"auth-api-go/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	config.ConnectDB()

	// Get PORT from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()

	// Health Check Route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Register API routes
	routes.AuthRoutes(r)  // ✅ Ensure this is here to register authentication routes

	fmt.Println("🚀 Server started on port", port)
	r.Run(":" + port)
}
