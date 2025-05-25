// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configurations.
type Config struct {
	MongoURI  string
	JWTSecret string
	Port      string
}

// LoadConfig reads configuration from .env file and environment variables.
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Make sure it exists and is accessible.")
	}

	// Retrieve environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set.")
	}
	if len(jwtSecret) < 32 { // Good practice for a strong secret
		log.Fatal("JWT_SECRET should be at least 32 characters long for security.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}

	return &Config{
		MongoURI:  mongoURI,
		JWTSecret: jwtSecret,
		Port:      port,
	}
}