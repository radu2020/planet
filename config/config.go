package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) ConnectionInfo() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "postgres",
		Port:     5432,
		User:     "user",
		Password: "password",
		Name:     "mydb",
	}
}

type Config struct {
	Port      int            `json:"port"`
	Env       string         `json:"env"`
	Database  PostgresConfig `json:"database"`
	FilePath  string         `json:"file_path"`
	BatchSize int            `json:"batch_size"` // Number of records per batch to be inserted in the db
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() Config {
	c := Config{
		Port:      getEnvInt("API_PORT", 8080),
		Env:       getEnv("ENV", "dev"),
		Database:  loadPostgresConfig(),
		FilePath:  getEnv("FILE_PATH", "/app/data/sample.csv"),
		BatchSize: getEnvInt("BATCH_SIZE", 50),
	}

	log.Println("Successfully loaded configuration.")
	return c
}

// loadPostgresConfig reads the PostgreSQL configuration from environment variables
func loadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     getEnv("POSTGRES_HOST", "postgres"),
		Port:     getEnvInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "user"),
		Password: getEnv("POSTGRES_PASSWORD", "password"),
		Name:     getEnv("POSTGRES_DB", "mydb"),
	}
}

// getEnv reads an environment variable or returns the default value if it's not set.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvInt reads an integer environment variable or returns the default value if it's not set.
func getEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}
