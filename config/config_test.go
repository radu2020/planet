package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Test ConnectionInfo format
func TestPostgresConfig_ConnectionInfo(t *testing.T) {
	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		Name:     "mydb",
	}

	expected := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
	assert.Equal(t, expected, cfg.ConnectionInfo())
}

// Test default PostgresConfig values
func TestDefaultPostgresConfig(t *testing.T) {
	cfg := DefaultPostgresConfig()
	assert.Equal(t, "postgres", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "user", cfg.User)
	assert.Equal(t, "password", cfg.Password)
	assert.Equal(t, "mydb", cfg.Name)
}

// Test IsProd function
func TestConfig_IsProd(t *testing.T) {
	devConfig := Config{Env: "dev"}
	prodConfig := Config{Env: "prod"}

	assert.False(t, devConfig.IsProd())
	assert.True(t, prodConfig.IsProd())
}

// Test LoadConfig with environment variables
func TestLoadConfig(t *testing.T) {
	os.Setenv("API_PORT", "9090")
	os.Setenv("ENV", "prod")
	os.Setenv("FILE_PATH", "/sample/data.csv")
	os.Setenv("BATCH_SIZE", "100")
	os.Setenv("POSTGRES_HOST", "db-host")
	os.Setenv("POSTGRES_PORT", "6543")
	os.Setenv("POSTGRES_USER", "admin")
	os.Setenv("POSTGRES_PASSWORD", "securepass")
	os.Setenv("POSTGRES_DB", "testdb")

	cfg := LoadConfig()

	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "prod", cfg.Env)
	assert.Equal(t, "/sample/data.csv", cfg.FilePath)
	assert.Equal(t, 100, cfg.BatchSize)
	assert.Equal(t, "db-host", cfg.Database.Host)
	assert.Equal(t, 6543, cfg.Database.Port)
	assert.Equal(t, "admin", cfg.Database.User)
	assert.Equal(t, "securepass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.Name)
}

// Test getEnv helper function
func TestGetEnv(t *testing.T) {
	os.Setenv("EXISTING_VAR", "value")
	assert.Equal(t, "value", getEnv("EXISTING_VAR", "default"))
	assert.Equal(t, "default", getEnv("MISSING_VAR", "default"))
}

// Test getEnvInt helper function// Test getEnvInt helper function
func TestGetEnvInt(t *testing.T) {
	os.Setenv("EXISTING_INT", "42")
	assert.Equal(t, 42, getEnvInt("EXISTING_INT", 10))
	assert.Equal(t, 10, getEnvInt("MISSING_INT", 10)) // Default value
	os.Setenv("INVALID_INT", "not_a_number")
	assert.Equal(t, 10, getEnvInt("INVALID_INT", 10)) // Should return default
}
