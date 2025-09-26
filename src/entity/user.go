package entity

import (
    "time"
    "gorm.io/gorm"
)

// User represents a user record in the database.
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `json:"name"`
    Email     string    `gorm:"unique" json:"email"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate sets timestamps before inserting a new record.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    now := time.Now().UTC()
    if u.CreatedAt.IsZero() {
        u.CreatedAt = now
    }
    u.UpdatedAt = now
    return nil
}

// BeforeUpdate updates the UpdatedAt timestamp before updating an existing record.
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    u.UpdatedAt = time.Now().UTC()
    return nil
}