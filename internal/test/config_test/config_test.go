package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Clear environment variables to test defaults
	os.Clearenv()

	// Set required DB environment variables
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "testuser")
	_ = os.Setenv("DB_PASSWORD", "testpass")

	cfg := config.Load()

	// Test App defaults
	assert.Equal(t, "Cocoon", cfg.App.Name)
	assert.Equal(t, "debug", cfg.Env)
	assert.Equal(t, "50051", cfg.App.Port)
	assert.Equal(t, 10*time.Second, cfg.Timeout)

	// Test Auth defaults
	assert.Equal(t, "thisissecret", cfg.SecretKey)
	assert.Equal(t, 24*time.Hour, cfg.TokenDuration)
	assert.Equal(t, "cocoon", cfg.Issuer)
	assert.Equal(t, 10, cfg.HashCost)

	// Test DB defaults
	assert.Equal(t, "postgres", cfg.Driver)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.DB.Port)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpass", cfg.Password)
	assert.Equal(t, "cocoon", cfg.DB.Name)
}

func TestConfig_CustomValues(t *testing.T) {
	// Set custom environment variables
	_ = os.Setenv("APP_NAME", "CustomApp")
	_ = os.Setenv("APP_ENV", "production")
	_ = os.Setenv("APP_PORT", "8080")
	_ = os.Setenv("APP_TIMEOUT", "30s")

	_ = os.Setenv("AUTH_SECRET_KEY", "customsecret")
	_ = os.Setenv("AUTH_TOKEN_DURATION", "48h")
	_ = os.Setenv("AUTH_ISSUER", "customissuer")
	_ = os.Setenv("AUTH_HASH_COST", "12")

	_ = os.Setenv("DB_DRIVER", "mysql")
	_ = os.Setenv("DB_HOST", "customhost")
	_ = os.Setenv("DB_PORT", "3306")
	_ = os.Setenv("DB_USER", "customuser")
	_ = os.Setenv("DB_PASSWORD", "custompass")
	_ = os.Setenv("DB_NAME", "customdb")

	cfg := config.Load()

	// Test App custom values
	assert.Equal(t, "CustomApp", cfg.App.Name)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, "8080", cfg.App.Port)
	assert.Equal(t, 30*time.Second, cfg.Timeout)

	// Test Auth custom values
	assert.Equal(t, "customsecret", cfg.SecretKey)
	assert.Equal(t, 48*time.Hour, cfg.TokenDuration)
	assert.Equal(t, "customissuer", cfg.Issuer)
	assert.Equal(t, 12, cfg.HashCost)

	// Test DB custom values
	assert.Equal(t, "mysql", cfg.Driver)
	assert.Equal(t, "customhost", cfg.Host)
	assert.Equal(t, "3306", cfg.DB.Port)
	assert.Equal(t, "customuser", cfg.User)
	assert.Equal(t, "custompass", cfg.Password)
	assert.Equal(t, "customdb", cfg.DB.Name)

	// Clean up
	os.Clearenv()
}
