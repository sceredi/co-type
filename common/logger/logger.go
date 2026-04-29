// Package logger contains logging utilities.
package logger

import (
	"log/slog"
	"os"
)

// Setup sets up the default logger for cli services.
func Setup() {
	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	logger := slog.New(
		slog.NewJSONHandler(os.Stderr, opts),
	)
	slog.SetDefault(logger)
}
