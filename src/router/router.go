package router

import (
    "github.com/gin-gonic/gin"
    "boilerplate-golang/src/controllers"
)

func InitRouter(r *gin.Engine) {
    // attach default middlewares and utility endpoints
    ApplyDefaults(r)
    api := r.Group("/api")

    api.GET("/users", controllers.GetUsers)
    api.POST("/users", controllers.CreateUser)
}
