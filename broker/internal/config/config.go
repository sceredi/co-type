// Package config is responsible for setting up the configuration for the broker application.
package config

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/joho/godotenv"
	handler "github.com/sceredi/co-type/broker/internal/grpc"
	"github.com/sceredi/co-type/broker/internal/repository/memory"
	"github.com/sceredi/co-type/broker/internal/service"
	"github.com/sceredi/co-type/common/logger"
	"github.com/sceredi/co-type/common/proto/control"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Setup initializes the configuration.
func Setup() error {
	logger.Setup()
	err := godotenv.Load()
	return err
}

// CreateDiscoveryListener creates and starts the gRPC server for the discovery service. It listens on the port specified in the environment variable "BROKER_PORT" and registers the DiscoveryServiceServer with the gRPC server.
func CreateDiscoveryListener(addr string, lis net.Listener) error {
	serverRepo := memory.NewServerRepository()
	serverSvc := service.NewServerService(serverRepo)
	controlHandler := handler.NewControlHandler(serverSvc)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    3 * time.Second,
			Timeout: 3 * time.Second,
		}),
	)
	control.RegisterControlServiceServer(grpcServer, controlHandler)

	reflection.Register(grpcServer)

	slog.InfoContext(context.Background(), "broker listening",
		slog.String("port", addr),
	)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
