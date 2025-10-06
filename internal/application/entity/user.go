package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user record in the database.
type User struct {
	ID        string         `json:"id" gorm:"column:id;primaryKey;type:varchar(255);comment:'Primary Key'"`
	Username  string         `json:"username" gorm:"column:username;type:varchar(50);comment:'username to login'"`
	Email     string         `json:"email" gorm:"column:email;type:varchar(100);comment:'email to login'"`
	Password  string         `json:"password" gorm:"column:password;type:varchar(255);comment:'password to login'"`
	FullName  string         `json:"full_name" gorm:"column:full_name;type:varchar(100);comment:'full name'"`
	IsActive  bool           `json:"is_active" gorm:"column:is_active;type:boolean;comment:'is active'"`
	IsAdmin   bool           `json:"is_admin" gorm:"column:is_admin;type:boolean;comment:'is admin'"`
	LastLogin *time.Time     `json:"last_login" gorm:"column:last_login;type:timestamp;comment:'last login'"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;comment:'created at'"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;comment:'updated at'"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp;comment:'deleted at'"`
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
	// Set default value for IsActive
	u.IsActive = true // Default to active
	return nil
}

// BeforeUpdate updates the UpdatedAt timestamp before updating an existing record.
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().UTC()
	return nil
}
