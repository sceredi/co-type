// Package app provides utilities for setting up the broker application
package app

import (
	"context"
	"log/slog"
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

const (
	heartbeatInterval = 10 * time.Second
	heartbeatTTL      = 3 * heartbeatInterval
)

// App represents the main application structure for the broker, containing everything needed to run the broker.
type App struct {
	Server    *grpc.Server
	cancelCtx context.CancelFunc
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

	ctx, cancel := context.WithCancel(context.Background())
	go runEviction(ctx, serverSvc)

	return &App{
		Server:    grpcServer,
		cancelCtx: cancel,
	}, nil
}

// Shutdown gracefully shuts down the application by stopping the gRPC server.
func (a *App) Shutdown() {
	a.cancelCtx()
	a.Server.GracefulStop()
}

func runEviction(ctx context.Context, serverSvc service.ServerService) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			n := serverSvc.EvictStale(heartbeatTTL)
			if n > 0 {
				slog.InfoContext(ctx, "Evicted stale servers", slog.Int("count", n))
			}
		}
	}
}
