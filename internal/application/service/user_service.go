package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"boilerplate-golang/internal/application/entity"
	"boilerplate-golang/internal/infrastructure/dbmanager"
)

type userService struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserService
func NewUserService() UserService {
	return &userService{
		db: dbmanager.DB(),
	}
}

// GetAllUsers returns a paginated list of users
func (s *userService) GetAllUsers(page, pageSize int) ([]entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var users []entity.User
	var total int64

	// Count total records
	err := s.db.Model(&entity.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = s.db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(id string) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	var user entity.User
	err := s.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user
func (s *userService) CreateUser(user entity.User) (*entity.User, error) {
	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// TODO: Hash password before saving

	err := s.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(user entity.User) (*entity.User, error) {
	if user.ID == 0 {
		return nil, errors.New("user ID cannot be empty")
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	err := s.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser deletes a user by their ID
func (s *userService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("user ID cannot be empty")
	}

	// Check if user exists
	_, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// Delete user
	return s.db.Delete(&entity.User{}, "id = ?", id).Error
}
