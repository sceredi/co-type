// Package main is the entry point of the application.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/config"
	"github.com/sceredi/co-type/client/internal/tui"
	cfg_utils "github.com/sceredi/co-type/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	f, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	defer func() {
		if err := f.Close(); err != nil {
			slog.InfoContext(context.Background(), "failed to close file",
				slog.String("err", err.Error()),
			)
		}
	}()

	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	logger := slog.New(
		slog.NewJSONHandler(f, opts),
	)
	slog.SetDefault(logger)

	cfg_utils.Setup()

	addr := os.Getenv("DISCOVERY_ADDR")
	port, err := strconv.Atoi(os.Getenv("DISCOVERY_PORT"))
	if err != nil {
		log.Fatalf("Invalid discovery port value: %v", err)
	}
	addr = fmt.Sprintf("%s:%d", addr, port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating discovery connection: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error closing gRPC connection: %v", err))
		}
	}()
	ds := config.CreateDiscoveryService(conn)

	m := tui.New(ds)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
