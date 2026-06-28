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
type DiscoveryService interface {
	// Register registers the server with the broker. lobbyIDs is empty on first startup and
	// contains the current lobby IDs when re-registering after a broker restart.
	Register(name, host string, port int, lobbyIDs []string) error

	// RegisterLobby registers a lobby with the service discovery mechanism using the provided lobby ID and server name.
	RegisterLobby(lobbyID, serverName string) error

	// StartHeartbeat starts a background goroutine that periodically sends heartbeats to the broker.
	// On any heartbeat failure it automatically re-registers the server (with current lobby IDs).
	// It stops when ctx is cancelled.
	StartHeartbeat(ctx context.Context, name, host string, port, load int, lobbyIDs func() []string)
}

type discoveryService struct {
	gtw gateway.ControlGateway
}

// NewDiscoveryService creates a new instance of DiscoveryService with the provided ControlGateway and returns it.
func NewDiscoveryService(gtw gateway.ControlGateway) DiscoveryService {
	return &discoveryService{gtw: gtw}
}

func (s *discoveryService) Register(name, host string, port int, lobbyIDs []string) error {
	return s.gtw.RegisterServer(name, host, port, lobbyIDs)
}

func (s *discoveryService) RegisterLobby(lobbyID, serverName string) error {
	return s.gtw.RegisterLobby(lobbyID, serverName)
}

func (s *discoveryService) StartHeartbeat(ctx context.Context, name, host string, port, load int, lobbyIDs func() []string) {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.gtw.Heartbeat(name, load); err != nil {
					slog.WarnContext(ctx, "Heartbeat failed, re-registering with broker",
						slog.String("name", name),
						slog.String("error", err.Error()),
					)
					if regErr := s.gtw.RegisterServer(name, host, port, lobbyIDs()); regErr != nil {
						slog.ErrorContext(ctx, "Re-registration failed",
							slog.String("name", name),
							slog.String("error", regErr.Error()),
						)
					}
				}
			}
		}
	}()
}
