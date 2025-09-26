package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user record in the database.
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // - means don't include in JSON
	FullName  string         `gorm:"size:100" json:"full_name,omitempty"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate sets timestamps and performs any pre-creation logic.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now().UTC()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	// Set default values
	if u.IsActive != false {
		u.IsActive = true // Default to active
	}
	return nil
}

// BeforeUpdate updates the UpdatedAt timestamp before updating an existing record.
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// Sanitize removes sensitive information from the user object
func (u *User) Sanitize() {
	u.Password = ""
}

// ToResponse converts User to UserResponse DTO
func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"full_name":  u.FullName,
		"is_active":  u.IsActive,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
