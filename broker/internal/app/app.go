// Package app provides utilities for setting up the broker application
package app

import (
	"time"

	"github.com/sceredi/co-type/broker/internal/api/grpc/handler"
	"github.com/sceredi/co-type/broker/internal/config"
	"github.com/sceredi/co-type/broker/internal/repository/memory"
	"github.com/sceredi/co-type/broker/internal/service"
	"github.com/sceredi/co-type/common/app"
	"github.com/sceredi/co-type/common/proto/control"
	"github.com/sceredi/co-type/common/proto/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// App represents the main application structure for the broker, containing everything needed to run the broker.
type App struct {
	Server *grpc.Server
}

// New initializes a new instance of App.
func New(_ *config.Config) (*App, error) {
	serverRepo := memory.NewServerRepository()
	serverSvc := service.NewServerService(serverRepo)
	lobbyRepo := memory.NewLobbyRepository()
	lobbySvc := service.NewLobbyService(lobbyRepo, serverSvc)

	controlHandler := handler.NewControlHandler(serverSvc, lobbySvc)
	discoveryHandler := handler.NewDiscoveryHandler(serverSvc, lobbySvc)

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

	return &App{
		Server: grpcServer,
	}, nil
}

// Shutdown gracefully shuts down the application by stopping the gRPC server.
func (a *App) Shutdown() {
	a.Server.GracefulStop()
}
