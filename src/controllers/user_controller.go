package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-mvcs-boilerplate/models"
	"go-mvcs-boilerplate/services"
	"go-mvcs-boilerplate/utils"
)

func GetUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	utils.SuccessResponse(c, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	newUser, err := services.CreateUser(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.SuccessResponse(c, newUser)
}
