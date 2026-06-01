// Package config is responsible for setting up the configuration for the server application.
package config

import (
	"context"
	"log"
	"time"

	"github.com/sceredi/co-type/common/config"
	"github.com/sceredi/co-type/common/proto/control"
	"github.com/sceredi/co-type/common/proto/lobby"
	"github.com/sceredi/co-type/server/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/server/internal/api/grpc/handler"
	"github.com/sceredi/co-type/server/internal/repository/memory"
	"github.com/sceredi/co-type/server/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// CreateDiscoveryService creates a new instance of DiscoveryService using the provided gRPC connection.
func CreateDiscoveryService(conn *grpc.ClientConn) service.DiscoveryService {
	c := control.NewControlServiceClient(conn)
	stream, err := c.Manage(context.Background())
	if err != nil {
		log.Fatalf("Error creating gRPC stream: %v", err)
	}
	gtw := gateway.NewControlGateway(stream)
	return service.NewDiscoveryService(gtw)
}

// CreateListeners creates and starts the gRPC servers.
func CreateListeners() *grpc.Server {
	lobbyRepo := memory.NewLobbyRepository()
	lobbySvc := service.NewLobbyService(lobbyRepo)

	lobbyHandler := handler.NewLobbyHandler(lobbySvc)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    3 * time.Second,
			Timeout: 3 * time.Second,
		}),
	)

	lobby.RegisterLobbyServiceServer(grpcServer, lobbyHandler)

	reflection.Register(grpcServer)

	config.CreateListener(grpcServer, "game")
	return grpcServer
}
