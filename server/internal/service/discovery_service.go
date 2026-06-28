// Package service contains the service layer of the server service.
// The service layer is responsible for implementing the business logic of the server service.
// It defines the interfaces that the handlers will use to interact with the domain models and repositories.
package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/sceredi/co-type/server/internal/api/grpc/gateway"
)

const heartbeatInterval = 10 * time.Second

// DiscoveryService defines the interface for managing service discovery in the server service.
// It provides a method to register the server with the service discovery mechanism.
// The Register method is responsible for registering the server and returns an error if the operation fails.
type DiscoveryService interface {
	// Register registers the server with the service discovery mechanism using the provided host and port.
	Register(name, host string, port int) error

	// RegisterLobby registers a lobby with the service discovery mechanism using the provided lobby ID and server name.
	RegisterLobby(lobbyID, serverName string) error

	// StartHeartbeat starts a background goroutine that periodically sends heartbeats to the broker.
	// It stops when ctx is cancelled.
	StartHeartbeat(ctx context.Context, name string, load int)
}

type discoveryService struct {
	gtw gateway.ControlGateway
}

// NewDiscoveryService creates a new instance of DiscoveryService with the provided ControlGateway and returns it.
func NewDiscoveryService(gtw gateway.ControlGateway) DiscoveryService {
	return &discoveryService{gtw: gtw}
}

func (s *discoveryService) Register(name, host string, port int) error {
	return s.gtw.RegisterServer(name, host, port)
}

func (s *discoveryService) RegisterLobby(lobbyID, serverName string) error {
	return s.gtw.RegisterLobby(lobbyID, serverName)
}

func (s *discoveryService) StartHeartbeat(ctx context.Context, name string, load int) {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.gtw.Heartbeat(name, load); err != nil {
					slog.WarnContext(ctx, "Heartbeat failed",
						slog.String("name", name),
						slog.String("error", err.Error()),
					)
				}
			}
		}
	}()
}
