package router

import (
	"github.com/gin-gonic/gin"

	"boilerplate-golang/internal/controller"
	"boilerplate-golang/internal/service"
)

// Register registers all HTTP routes on the given engine.
func Register(r *gin.Engine) {
	// Initialize services
	userService := service.NewUserService()

	// Initialize controllers
	userController := controller.NewUserController(userService)

	// API routes
	api := r.Group("/api")
	{
		// User routes
		userController.RegisterRoutes(api)

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}
