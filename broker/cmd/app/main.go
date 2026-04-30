// Package main is the entrypoint of the broker application.
package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sceredi/co-type/broker/internal/config"
	cfg_utils "github.com/sceredi/co-type/common/config"
)

func main() {
	writer := io.MultiWriter(os.Stderr)
	if os.Getenv("LOCAL") == "true" {
		f, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("failed to close file: %v", err)
			}
		}()

		writer = io.MultiWriter(os.Stderr, f)
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	logger := slog.New(
		slog.NewJSONHandler(writer, opts),
	)
	slog.SetDefault(logger)

	cfg_utils.Setup()

	grpcServer := config.CreateListeners()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()
}
