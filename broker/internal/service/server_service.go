// Package service contains the service layer of the broker service.
// The service layer is responsible for implementing the business logic of the broker service.
// It defines the interfaces that the handlers will use to interact with the domain models and repositories.
package service

import (
	"context"
	"log/slog"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

// ServerService defines the interface for managing servers in the broker service.
type ServerService interface {
	// Create creates a new server in the broker service. It takes a CreateServerRequest and returns the created Server or an error if the operation fails.
	Create(name, addr string, port int) (*domain.Server, error)
	// LowestLoad finds the server with the lowest load. It returns it or an error if there are no available servers.
	LowestLoad() (*domain.Server, error)
}

type serverService struct {
	serverRepo repository.ServerRepository
}

// NewServerService creates a new instance of ServerService with the provided ServerRepository and returns it.
func NewServerService(serverRepo repository.ServerRepository) ServerService {
	return &serverService{serverRepo: serverRepo}
}

func (s *serverService) Create(name, addr string, port int) (*domain.Server, error) {
	slog.DebugContext(context.Background(), "Creating server",
		slog.String("name", name),
		slog.String("addr", addr),
		slog.Int("port", port),
	)
	server, err := domain.NewServer(name, addr, port)
	if err != nil {
		return nil, err
	}
	return s.serverRepo.Create(server)
}

func (s *serverService) LowestLoad() (*domain.Server, error) {
	slog.DebugContext(context.Background(), "Finding server with lowest load")
	svs := s.serverRepo.List()
	if len(svs) == 0 {
		return nil, domain.ErrNoAvailableServers
	}
	lowest := svs[0]
	for _, server := range svs {
		if server.Load < lowest.Load {
			lowest = server
		}
	}
	return lowest, nil
}
