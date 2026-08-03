package memory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

type lobbyRepository struct {
	mu      sync.RWMutex
	lobbies map[repository.LobbyID]repository.ServerName
}

// NewLobbyRepository creates a new instance of LobbyRepository and returns it.
func NewLobbyRepository() repository.LobbyRepository {
	return &lobbyRepository{
		lobbies: make(map[repository.LobbyID]repository.ServerName),
	}
}

func (r *lobbyRepository) Create(lobbyID repository.LobbyID, serverName repository.ServerName) error {
	slog.DebugContext(context.Background(), "Storing lobby in memory",
		slog.String("lobby_id", string(lobbyID)),
		slog.String("server_name", string(serverName)),
	)
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.lobbies[lobbyID]
	if ok {
		return domain.ErrLobbyAlreadyExists
	}
	r.lobbies[lobbyID] = serverName
	return nil
}

func (r *lobbyRepository) Upsert(lobbyID repository.LobbyID, serverName repository.ServerName) {
	slog.DebugContext(context.Background(), "Upserting lobby in memory",
		slog.String("lobby_id", string(lobbyID)),
		slog.String("server_name", string(serverName)),
	)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lobbies[lobbyID] = serverName
}

func (r *lobbyRepository) Get(id repository.LobbyID) (repository.ServerName, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	serverName, ok := r.lobbies[id]
	if !ok {
		return "", domain.ErrLobbyNotFound
	}
	return serverName, nil
}

func (r *lobbyRepository) Delete(id repository.LobbyID) error {
	slog.DebugContext(context.Background(), "Deleting lobby from memory",
		slog.String("lobby_id", string(id)),
	)
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.lobbies[id]
	if !ok {
		return domain.ErrLobbyNotFound
	}
	delete(r.lobbies, id)
	return nil
}
