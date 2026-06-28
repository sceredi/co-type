// Package app provides utilities for setting up the server application
package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sceredi/co-type/common/app"
	"github.com/sceredi/co-type/common/proto/control"
	"github.com/sceredi/co-type/common/proto/lobby"
	"github.com/sceredi/co-type/server/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/server/internal/api/grpc/handler"
	"github.com/sceredi/co-type/server/internal/config"
	"github.com/sceredi/co-type/server/internal/repository/memory"
	"github.com/sceredi/co-type/server/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// App represents the main application structure for the server, containing everything needed to run the server.
type App struct {
	Server *grpc.Server
	Conn   *grpc.ClientConn
	// TODO: wrap into a "grpc.client.{discovery_client, connection_supervisor}" that does healthchecks and reconnection
	DiscoveryService service.DiscoveryService
	cancelCtx        context.CancelFunc
}

// New initializes a new instance of App.
func New(cfg *config.Config) (*App, error) {
	conn, err := createConn(cfg)
	if err != nil {
		return nil, err
	}
	c := control.NewControlServiceClient(conn)

	lobbyRepo := memory.NewLobbyRepository()

	ctx, cancel := context.WithCancel(context.Background())
	ctrlGtw := gateway.NewControlGateway(ctx, c)

	discoverySvc := service.NewDiscoveryService(ctrlGtw)
	lobbySvc := service.NewLobbyService(cfg.Name, ctrlGtw, lobbyRepo)

	lobbyHandler := handler.NewLobbyHandler(lobbySvc)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    3 * time.Second,
			Timeout: 3 * time.Second,
		}),
	)

	lobby.RegisterLobbyServiceServer(grpcServer, lobbyHandler)

	reflection.Register(grpcServer)

	app.CreateListener(grpcServer, "game")

	// TODO: same as todo above
	err = discoverySvc.Register(cfg.Name, cfg.Addr, cfg.Port)
	if err != nil {
		cancel()
		if closeErr := conn.Close(); closeErr != nil {
			slog.Error(fmt.Sprintf("Error closing gRPC connection after registration failure: %v", closeErr))
		}
		return nil, err
	}

	discoverySvc.StartHeartbeat(ctx, cfg.Name, 0)

	return &App{
		Server:           grpcServer,
		Conn:             conn,
		DiscoveryService: discoverySvc,
		cancelCtx:        cancel,
	}, nil
}

// Shutdown gracefully shuts down the application by stopping the gRPC server and closing the connection.
func (a *App) Shutdown() {
	a.cancelCtx()
	a.Server.GracefulStop()
	if err := a.Conn.Close(); err != nil {
		slog.Error(fmt.Sprintf("Error closing gRPC connection: %v", err))
	}
}

func createConn(cfg *config.Config) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.ControlAddr, cfg.ControlPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
