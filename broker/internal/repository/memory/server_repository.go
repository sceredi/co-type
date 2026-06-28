// Package memory provides an in-memory implementation of the repository interfaces defined in the repository package.
package memory

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

type serverEntry struct {
	server   *domain.Server
	lastSeen time.Time
}

type serverRepository struct {
	mu      sync.RWMutex
	entries []*serverEntry
}

// NewServerRepository creates a new instance of ServerRepository.
func NewServerRepository() repository.ServerRepository {
	return &serverRepository{
		entries: make([]*serverEntry, 0),
	}
}

func (r *serverRepository) Create(server *domain.Server) (*domain.Server, error) {
	slog.InfoContext(context.Background(), "Storing server in memory",
		slog.String("server", server.Addr),
		slog.Int("port", server.Port),
	)
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.entries {
		if e.server.Name == server.Name {
			return nil, domain.ErrServerAlreadyExists
		}
	}
	r.entries = append(r.entries, &serverEntry{server: server, lastSeen: time.Now()})
	return server, nil
}

func (r *serverRepository) List() []*domain.Server {
	r.mu.RLock()
	defer r.mu.RUnlock()
	servers := make([]*domain.Server, len(r.entries))
	for i, e := range r.entries {
		servers[i] = e.server
	}
	return servers
}

func (r *serverRepository) UpdateLastSeen(name string, load int, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.entries {
		if e.server.Name == name {
			e.server.Load = load
			e.lastSeen = t
			return nil
		}
	}
	return domain.ErrServerNotFound
}

func (r *serverRepository) Delete(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, e := range r.entries {
		if e.server.Name == name {
			r.entries = append(r.entries[:i], r.entries[i+1:]...)
			return nil
		}
	}
	return domain.ErrServerNotFound
}

// ListStale returns the names of servers whose lastSeen is older than the given TTL.
func (r *serverRepository) ListStale(ttl time.Duration) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var stale []string
	cutoff := time.Now().Add(-ttl)
	for _, e := range r.entries {
		if e.lastSeen.Before(cutoff) {
			stale = append(stale, e.server.Name)
		}
	}
	return stale
}
