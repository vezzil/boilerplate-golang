package routes

import (
	"github.com/gin-gonic/gin"
	"go-mvcs-boilerplate/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/users", controllers.GetUsers)
		api.POST("/users", controllers.CreateUser)
	}
}
