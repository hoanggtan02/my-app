package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
}

// LoadConfig reads configuration from .env file.
func LoadConfig(path string) (*Config, error) {
	// Construct the full path to the .env file
	envPath := path + "/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found at %s. Using environment variables.", envPath)
	}

	cfg := &Config{
		ServerPort:  os.Getenv("SERVER_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080" // Default port
	}

	return cfg, nil
}
