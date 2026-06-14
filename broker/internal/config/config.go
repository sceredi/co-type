// Package config is responsible for setting up the configuration for the broker application.
package config

import (
	"time"

	"github.com/sceredi/co-type/broker/internal/api/grpc/handler"
	"github.com/sceredi/co-type/broker/internal/repository/memory"
	"github.com/sceredi/co-type/broker/internal/service"
	"github.com/sceredi/co-type/common/app"
	"github.com/sceredi/co-type/common/proto/control"
	"github.com/sceredi/co-type/common/proto/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// CreateListeners creates and starts the gRPC servers.
func CreateListeners() *grpc.Server {
	serverRepo := memory.NewServerRepository()
	serverSvc := service.NewServerService(serverRepo)
	lobbyRepo := memory.NewLobbyRepository()
	lobbySvc := service.NewLobbyService(lobbyRepo, serverSvc)

	controlHandler := handler.NewControlHandler(serverSvc, lobbySvc)
	discoveryHandler := handler.NewDiscoveryHandler(serverSvc)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    3 * time.Second,
			Timeout: 3 * time.Second,
		}),
	)

	control.RegisterControlServiceServer(grpcServer, controlHandler)
	discovery.RegisterDiscoveryServiceServer(grpcServer, discoveryHandler)

	reflection.Register(grpcServer)

	app.CreateListener(grpcServer, "control")
	app.CreateListener(grpcServer, "discovery")
	return grpcServer
}
