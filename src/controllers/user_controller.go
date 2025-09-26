package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"boilerplate-golang/src/entity"
	"boilerplate-golang/src/services"
	"boilerplate-golang/src/tools"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		tools.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	tools.SuccessResponse(c, users)
}

func CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		tools.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	newUser, err := services.CreateUser(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SuccessResponse(c, newUser)
}
