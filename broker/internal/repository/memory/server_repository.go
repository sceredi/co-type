// Package memory provides an in-memory implementation of the repository interfaces defined in the repository package.
package memory

import (
	"context"
	"log/slog"

	"github.com/sceredi/co-type/broker/internal/domain"
	"github.com/sceredi/co-type/broker/internal/repository"
)

type serverRepository struct {
	servers []*domain.Server
}

// NewServerRepository creates a new instance of ServerRepository.
func NewServerRepository() repository.ServerRepository {
	return &serverRepository{
		servers: make([]*domain.Server, 0),
	}
}

func (r *serverRepository) Create(server *domain.Server) (*domain.Server, error) {
	slog.InfoContext(context.Background(), "Storing server in memory",
		slog.String("server", server.Addr),
		slog.Int("port", int(server.Port)),
	)
	for _, s := range r.servers {
		if s.Addr == server.Addr {
			return nil, domain.ErrServerAlreadyExists
		}
	}
	r.servers = append(r.servers, server)
	return server, nil
}
