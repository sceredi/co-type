// Package logger contains logging utilities.
package logger

import (
	"log/slog"
	"os"
)

// Setup sets up the default logger for cli services.
func Setup() {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stderr, opts),
	)
	slog.SetDefault(logger)
}
