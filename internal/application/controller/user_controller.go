package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"boilerplate-golang/internal/application/dto"
	"boilerplate-golang/internal/application/entity"
	"boilerplate-golang/internal/application/service"
	"boilerplate-golang/pkg/response"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService service.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// RegisterRoutes registers user routes
func (uc *UserController) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("", uc.GetUsers)
		userGroup.POST("", uc.CreateUser)
		userGroup.GET("/:id", uc.GetUser)
		userGroup.PUT("/:id", uc.UpdateUser)
		userGroup.DELETE("/:id", uc.DeleteUser)
	}
}

// GetUsers handles GET /api/users
func (uc *UserController) GetUsers(c *gin.Context) {
	// Parse pagination query params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Get users with pagination
	users, total, err := uc.userService.GetAllUsers(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}

	// Convert to DTOs using entity's ToResponse method
	userDTOs := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToResponse()
	}

	// Return paginated response
	response.Success(c, dto.PaginatedResponse{
		Success:  true,
		Data:     userDTOs,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

// GetUser handles GET /api/users/:id
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			response.NotFound(c, "User")
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		}
		return
	}

	response.Success(c, user.ToResponse())
}

// CreateUser handles POST /api/users
func (uc *UserController) CreateUser(c *gin.Context){
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input data")
		return
	}

	// Map DTO to entity
	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // In a real app, this should be hashed
		FullName: req.FullName,
	}

	createdUser, err := uc.userService.CreateUser(user)
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint" {
			response.Error(c, http.StatusConflict, "User with this email or username already exists")
		} else {
			response.InternalServerError(c, "Failed to create user: "+err.Error())
		}
		return
	}

	response.SuccessWithStatus(c, http.StatusCreated, createdUser.ToResponse(), "User created successfully")
}

// UpdateUser handles PUT /api/users/:id
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required")
		return
	}

	// Get existing user
	existingUser, err := uc.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			response.NotFound(c, "User")
		} else {
			response.InternalServerError(c, "Failed to fetch user: "+err.Error())
		}
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input data")
		return
	}

	// Update fields if provided in request
	if req.Username != nil {
		existingUser.Username = *req.Username
	}
	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.FullName != nil {
		existingUser.FullName = *req.FullName
	}

	updatedUser, err := uc.userService.UpdateUser(*existingUser)
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint" {
			response.Error(c, http.StatusConflict, "User with this email or username already exists")
		} else {
			response.InternalServerError(c, "Failed to update user: "+err.Error())
		}
		return
	}

	response.Success(c, updatedUser.ToResponse())
}

// DeleteUser handles DELETE /api/users/:id
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "User ID is required")
		return
	}

	err := uc.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "record not found" {
			response.NotFound(c, "User")
		} else {
			response.InternalServerError(c, "Failed to delete user: "+err.Error())
		}
		return
	}

	response.Success(c, nil)
}
