// Package main is the entrypoint of the broker application.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"

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
	controlPort, err := strconv.Atoi(os.Getenv("CONTROL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	addr := fmt.Sprintf(":%d", controlPort)
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", controlPort, err)
	}
	err = config.CreateDiscoveryListener(addr, lis)
	if err != nil {
		log.Fatalf("Failed to start discovery listener: %v", err)
	}
}
