package config

import (
	"fmt"
	"os"

	"log"

	"github.com/joho/godotenv"
)

// Config holds the essential configuration values
type Config struct {
	// Database
	PGHostname string
	PGPort     string
	PGDatabase string
	PGUsername string
	PGPassword string

	// Application
	AppEnv string
	Port   string
}

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Please create a .env file with the required environment variables.")
	}
}

// Load creates a new Config instance with values from environment variables
func Load() *Config {
	LoadEnv()

	return &Config{
		// Database
		PGHostname: getEnv("PG_HOSTNAME", "localhost"),
		PGPort:     getEnv("PG_PORT", "5432"),
		PGDatabase: getEnv("PG_DATABASE", "athenai"),
		PGUsername: getEnv("PG_USERNAME", "postgres"),
		PGPassword: getEnv("PG_PASSWORD", ""),

		// Application
		AppEnv: getEnv("APP_ENV", "prod"),
		Port:   getEnv("PORT", "8080"),
	}
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.PGHostname, c.PGPort, c.PGUsername, c.PGPassword, c.PGDatabase)
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "dev" || c.AppEnv == "development"
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "prod" || c.AppEnv == "production"
}

// Helper function
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
