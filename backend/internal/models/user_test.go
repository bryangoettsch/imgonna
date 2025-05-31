package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	return db
}

func TestUser_BeforeCreate(t *testing.T) {
	db := setupTestDB(t)

	user := &User{
		Auth0ID: "auth0|123456",
		Email:   "test@example.com",
		Name:    "Test User",
	}

	err := db.Create(user).Error
	assert.NoError(t, err)
	assert.Equal(t, RoleUser, user.Role)
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{"User role", RoleUser, false},
		{"Admin role", RoleAdmin, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.expected, user.IsAdmin())
		})
	}
}

func TestUser_UpdateLastLogin(t *testing.T) {
	db := setupTestDB(t)

	user := &User{
		Auth0ID:    "auth0|123456",
		Email:      "test@example.com",
		Name:       "Test User",
		LoginCount: 0,
	}

	err := db.Create(user).Error
	assert.NoError(t, err)

	err = user.UpdateLastLogin(db)
	assert.NoError(t, err)

	// Reload user from database
	var updatedUser User
	err = db.First(&updatedUser, user.ID).Error
	assert.NoError(t, err)

	assert.NotNil(t, updatedUser.LastLoginAt)
	assert.Equal(t, 1, updatedUser.LoginCount)
	assert.WithinDuration(t, time.Now(), *updatedUser.LastLoginAt, time.Minute)
}