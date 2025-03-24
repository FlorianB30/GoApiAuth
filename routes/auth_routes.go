package routes

import (
	"auth-api-go/controllers"
	"auth-api-go/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	group := r.Group("/633aa0c9-d914-4308-8fde-4b9333516586")
	{
		group.POST("/register", controllers.Register)
		group.POST("/login", controllers.Login)
		group.POST("/forgot-password", controllers.ForgotPassword)
		group.POST("/reset-password", controllers.ResetPassword)
		// ✅ Secure routes with authentication
		protected := group.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // Require authentication
		protected.GET("/me", controllers.GetUserProfile)
		protected.POST("/logout", controllers.Logout)
	}
}
