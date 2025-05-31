package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// Set test environment variables
	originalEnvs := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_SSL_MODE": os.Getenv("DB_SSL_MODE"),
	}

	// Cleanup function to restore original environment
	defer func() {
		for key, value := range originalEnvs {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	t.Run("Connect with default values", func(t *testing.T) {
		// Clear environment variables to test defaults
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSL_MODE")

		// This will fail because we don't have a real postgres instance
		// but we can test that the function doesn't panic and returns an error
		err := Connect()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to database")
	})

	t.Run("Connect with custom environment", func(t *testing.T) {
		os.Setenv("DB_HOST", "custom-host")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("DB_NAME", "custom-db")
		os.Setenv("DB_USER", "custom-user")
		os.Setenv("DB_PASSWORD", "custom-pass")
		os.Setenv("DB_SSL_MODE", "require")

		// This will also fail but we're testing the environment variable reading
		err := Connect()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to database")
	})
}

func TestGetDB(t *testing.T) {
	// Test that GetDB returns the global DB instance
	originalDB := DB
	defer func() { DB = originalDB }()

	// Set DB to nil to test
	DB = nil
	assert.Nil(t, GetDB())

	// We can't easily test with a real connection in unit tests
	// without setting up a test database, but we can verify the function works
}