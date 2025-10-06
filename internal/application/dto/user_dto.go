package dto

import (
	"boilerplate-golang/internal/application/entity"
	"time"
)

// // UserCreateRequest represents the request body for creating a new user
// type UserCreateRequest struct {
// 	Username string `json:"username" binding:"required,min=3,max=50"`
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=8"`
// 	FullName string `json:"full_name" binding:"required,min=2,max=100"`
// }

// // UserUpdateRequest represents the request body for updating a user
// type UserUpdateRequest struct {
// 	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
// 	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
// 	FullName *string `json:"full_name,omitempty" binding:"omitempty,min=2,max=100"`
// }

// UserResponse represents the user data sent in the response
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetUserResponse(entity entity.User) UserResponse {
	return UserResponse{
		ID:        entity.ID,
		Username:  entity.Username,
		Email:     entity.Email,
		FullName:  entity.FullName,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// // LoginRequest represents the login request payload
// type LoginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required"`
// }

// // TokenResponse represents the authentication token response
// type TokenResponse struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token,omitempty"`
// 	ExpiresIn    int64  `json:"expires_in"`
// 	TokenType    string `json:"token_type"`
// }
