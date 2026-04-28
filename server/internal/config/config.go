// Package config is responsible for setting up the configuration for the server application.
package config

import (
	"context"
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/sceredi/co-type/common/logger"
	"github.com/sceredi/co-type/common/proto/control"
	ctgrpc "github.com/sceredi/co-type/server/internal/grpc"
	"github.com/sceredi/co-type/server/internal/service"
	"google.golang.org/grpc"
)

// Setup initializes the configuration for the server application.
func Setup() {
	logger.Setup()
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}
}

// CreateDiscoveryService creates a new instance of DiscoveryService using the provided gRPC connection.
func CreateDiscoveryService(conn *grpc.ClientConn) service.DiscoveryService {
	c := control.NewControlServiceClient(conn)
	stream, err := c.Manage(context.Background())
	if err != nil {
		log.Fatalf("Error creating gRPC stream: %v", err)
	}
	gtw := ctgrpc.NewControlGateway(stream)
	return service.NewDiscoveryService(gtw)
}
