// Package main is the entry point of the server application.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	cfg_utils "github.com/sceredi/co-type/common/config"
	"github.com/sceredi/co-type/server/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	serverName := os.Getenv("SERVER_NAME")
	serverAddr := os.Getenv("SERVER_ADDR")
	serverPortStr := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		log.Fatalf("Error parsing server port: %v", err)
	}
	idx := serverName[len(serverName)-1] - '0'
	serverPort = serverPort + int(idx)
	slog.InfoContext(context.Background(), "Server info",
		slog.String("serverName", serverName),
		slog.String("serverAddr", serverAddr),
		slog.Int("idx", int(idx)),
		slog.Int("serverPort", serverPort),
	)
	controlAddr := os.Getenv("CONTROL_ADDR")
	controlPort := os.Getenv("CONTROL_PORT")
	addr := fmt.Sprintf("%s:%s", controlAddr, controlPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating gRPC client: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error closing gRPC connection: %v", err))
		}
	}()
	ds := config.CreateDiscoveryService(conn)
	err = ds.Register(serverName, serverAddr, serverPort)
	if err != nil {
		log.Fatalf("Error registering server: %v", err)
	}

	grpcServer := config.CreateListeners()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.GracefulStop()
}
