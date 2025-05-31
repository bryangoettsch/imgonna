package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Auth0 fields
	Auth0ID string `json:"auth0_id" gorm:"uniqueIndex;not null"`
	Email   string `json:"email" gorm:"uniqueIndex;not null"`
	Name    string `json:"name" gorm:"not null"`

	// Application fields
	Role   Role   `json:"role" gorm:"type:varchar(20);default:'user'"`
	Active bool   `json:"active" gorm:"default:true"`
	Avatar string `json:"avatar,omitempty"`

	// Metadata
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LoginCount  int        `json:"login_count" gorm:"default:0"`
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RoleUser
	}
	return nil
}

// IsAdmin returns true if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// UpdateLastLogin updates the last login timestamp and increments login count
func (u *User) UpdateLastLogin(tx *gorm.DB) error {
	now := time.Now()
	return tx.Model(u).Updates(map[string]interface{}{
		"last_login_at": &now,
		"login_count":   gorm.Expr("login_count + 1"),
	}).Error
}