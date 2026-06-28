// Package service contains the service layer of the broker service.
// The service layer is responsible for implementing the business logic of the broker service.
// It defines the interfaces that the handlers will use to interact with the domain models and repositories.
package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

// ServerService defines the interface for managing servers in the broker service.
type ServerService interface {
	// Create creates a new server in the broker service. It takes a CreateServerRequest and returns the created Server or an error if the operation fails.
	Create(name, addr string, port int) (*domain.Server, error)
	// Upsert registers or updates a server entry. It never returns ErrServerAlreadyExists.
	Upsert(name, addr string, port int) *domain.Server
	// LowestLoad finds the server with the lowest load. It returns it or an error if there are no available servers.
	LowestLoad() (*domain.Server, error)
	// LowestLoadExcluding finds the server with the lowest load, skipping the named server.
	// Falls back to LowestLoad if the excluded server is the only one registered.
	LowestLoadExcluding(excludeName string) (*domain.Server, error)
	// GetByName retrieves a server by its name. It returns the server or an error if the server is not found or if the operation fails.
	GetByName(name string) (*domain.Server, error)
	// Heartbeat updates the last-seen timestamp and load for a registered server.
	Heartbeat(name string, load int) error
	// EvictStale removes servers that have not sent a heartbeat within ttl and returns the count of evicted servers.
	EvictStale(ttl time.Duration) int
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

func (s *serverService) LowestLoadExcluding(excludeName string) (*domain.Server, error) {
	slog.DebugContext(context.Background(), "Finding server with lowest load excluding", slog.String("exclude", excludeName))
	svs := s.serverRepo.List()
	candidates := make([]*domain.Server, 0, len(svs))
	for _, sv := range svs {
		if sv.Name != excludeName {
			candidates = append(candidates, sv)
		}
	}
	if len(candidates) == 0 {
		// Only the crashed server is known — fall back to it.
		return s.LowestLoad()
	}
	lowest := candidates[0]
	for _, server := range candidates {
		if server.Load < lowest.Load {
			lowest = server
		}
	}
	return lowest, nil
}

func (s *serverService) GetByName(name string) (*domain.Server, error) {
	slog.DebugContext(context.Background(), "Retrieving server by name",
		slog.String("name", name),
	)
	svs := s.serverRepo.List()
	for _, server := range svs {
		if server.Name == name {
			return server, nil
		}
	}
	return nil, domain.ErrServerNotFound
}

func (s *serverService) Upsert(name, addr string, port int) *domain.Server {
	server := &domain.Server{Name: name, Addr: addr, Port: port}
	return s.serverRepo.Upsert(server)
}

func (s *serverService) Heartbeat(name string, load int) error {
	slog.DebugContext(context.Background(), "Received heartbeat",
		slog.String("name", name),
		slog.Int("load", load),
	)
	return s.serverRepo.UpdateLastSeen(name, load, time.Now())
}

func (s *serverService) EvictStale(ttl time.Duration) int {
	stale := s.serverRepo.ListStale(ttl)
	for _, name := range stale {
		if err := s.serverRepo.Delete(name); err != nil {
			slog.WarnContext(context.Background(), "Failed to evict stale server",
				slog.String("name", name),
				slog.String("error", err.Error()),
			)
			continue
		}
		slog.InfoContext(context.Background(), "Evicted stale server",
			slog.String("name", name),
		)
	}
	return len(stale)
}
