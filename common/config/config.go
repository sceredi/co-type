// Package config contains useful configuration utilities.
package config

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/sceredi/co-type/common/logger"
)

// Setup initializes the configuration.
func Setup() {
	logger.Setup()
	err := godotenv.Load()
	if err != nil {
		slog.Debug("No .env file found, using environment variables")
	}
}
