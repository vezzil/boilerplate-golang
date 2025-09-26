package service

import (
	"boilerplate-golang/internal/entity"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	// Get all users with pagination
	GetAllUsers(page, pageSize int) ([]entity.User, int64, error)
	// Get user by ID
	GetUserByID(id string) (*entity.User, error)
	// Create a new user
	CreateUser(user entity.User) (*entity.User, error)
	// Update an existing user
	UpdateUser(user entity.User) (*entity.User, error)
	// Delete a user
	DeleteUser(id string) error
}
