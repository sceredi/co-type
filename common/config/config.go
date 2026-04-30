// Package config contains useful configuration utilities.
package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

// Setup initializes the configuration.
func Setup() {
	err := godotenv.Load()
	if err != nil {
		slog.Debug("No .env file found, using environment variables")
	}
}
