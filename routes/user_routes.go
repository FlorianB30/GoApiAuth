package routes

import (
	"auth-api-go/controllers"
	"auth-api-go/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	group := r.Group("/633aa0c9-d914-4308-8fde-4b9333516586")
	{
		protected := group.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // Require authentication
		protected.GET("/me", controllers.GetUserProfile)
	}
}
