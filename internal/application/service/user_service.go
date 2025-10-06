package service

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"boilerplate-golang/internal/application/dto"
	"boilerplate-golang/internal/application/entity"
	"boilerplate-golang/internal/application/tools"
	"boilerplate-golang/internal/infrastructure/dbmanager"
	"boilerplate-golang/internal/infrastructure/logger"
)

type userService struct {
}

// GetAllUsers returns a paginated list of users with filtering and sorting
func (s *userService) GetAllUsers(page, pageSize int, filterSearch, sortBy, sortOrder string) dto.ResponseDto {
	var users []entity.User
	var totalCount int64

	db := dbmanager.GetDB()

	// Base query
	query := db.Model(&entity.User{}).Where("deleted_at IS NULL")

	// Apply search filter
	if filterSearch != "" {
		searchTerm := "%" + filterSearch + "%"
		query = query.Where(
			"LOWER(username) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(full_name) LIKE LOWER(?)",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// Apply sorting
	if sortBy != "" && (sortOrder == "asc" || sortOrder == "desc") {
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		// Default sorting
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if pageSize > 0 {
		offset := (page - 1) * pageSize
		if err := query.Count(&totalCount).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
			logger.Error("Error fetching users with pagination: %v", err)
			return *dto.Fail("Error fetching users with pagination")
		}
	} else {
		// If no pagination, fetch all records
		if err := query.Count(&totalCount).Find(&users).Error; err != nil {
			logger.Error("Error fetching users with no pagination: %v", err)
			return *dto.Fail("Error fetching users with no pagination")
		}
	}

	// Convert entities to DTOs
	userDtos := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userDtos[i] = dto.GetUserResponse(user)
	}

	// Return paginated response
	return *dto.SuccessCount(userDtos, totalCount)
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(id string) dto.ResponseDto {
	var user entity.User

	db := dbmanager.GetDB()
	if err := db.First(&user, id).Error; err != nil {
		logger.Error("Error fetching user by ID: %v", err)
		return *dto.Fail("Error fetching user by ID")
	}

	return *dto.Success(dto.GetUserResponse(user))
}

// CreateUser creates a new user
func (s *userService) CreateUser(username, email, password, fullName string) dto.ResponseDto {
	db := dbmanager.GetDB()
	tx := db.Begin()

	// Validate required fields
	if username == "" || email == "" || password == "" || fullName == "" {
		return *dto.Fail("All fields are required")
	}

	// Check if username contains spaces
	if strings.Contains(username, " ") {
		return *dto.Fail("Username should not contain spaces")
	}

	// Check if username already exists
	var existingUser entity.User
	if err := tx.Where("username = ?", username).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return *dto.Fail("Username already exists")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		logger.Error("Error checking username existence: %v", err)
		return *dto.Fail("Error checking username availability")
	}

	// Validate email format
	if !tools.IsValidEmail(email) {
		return *dto.Fail("Invalid email format")
	}

	// Check if email already exists
	if err := tx.Where("email = ?", email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return *dto.Fail("Email already in use")
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		logger.Error("Error checking email existence: %v", err)
		return *dto.Fail("Error checking email availability")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		logger.Error("Error hashing password: %v", err)
		return *dto.Fail("Error creating user account")
	}

	// Create new user
	newUser := entity.User{
		ID:       tools.NewUuid(),
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		FullName: fullName,
		IsActive: true,
		IsAdmin:  false,
	}

	// Save user to database
	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		logger.Error("Error creating user: %v", err)
		return *dto.Fail("Error creating user account")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error("Error committing transaction: %v", err)
		return *dto.Fail("Error creating user account")
	}

	// Return the created user (without sensitive data)
	return *dto.Success(dto.GetUserResponse(newUser))
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(id, username, email, password, fullName string) dto.ResponseDto {
	db := dbmanager.GetDB()

	var user entity.User
	if err := db.First(&user, id).Error; err != nil {
		return *dto.Fail("User not found")
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		user.Password = password
	}
	if fullName != "" {
		user.FullName = fullName
	}

	if err := db.Save(&user).Error; err != nil {
		logger.Error("Error updating user: %v", err)
		return *dto.Fail("Error updating user")
	}

	return *dto.Success("User updated successfully")
}

// SoftDeleteUser soft deletes a user by their ID
func (s *userService) SoftDeleteUser(id string) dto.ResponseDto {
	db := dbmanager.GetDB()

	var user entity.User
	if err := db.First(&user, id).Error; err != nil {
		logger.Error("Error soft deleting user: %v", err)
		return *dto.Fail("User not found")
	}

	user.IsActive = false
	user.DeletedAt = gorm.DeletedAt{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	if err := db.Save(&user).Error; err != nil {
		logger.Error("Error soft deleting user: %v", err)
		return *dto.Fail("Error soft deleting user")
	}

	return *dto.Success("User soft deleted successfully")
}
