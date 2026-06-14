package memory

import (
	"context"
	"log/slog"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

type lobbyRepository struct {
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
	_, ok := r.lobbies[lobbyID]
	if ok {
		return domain.ErrLobbyAlreadyExists
	}
	r.lobbies[lobbyID] = serverName
	return nil
}

func (r *lobbyRepository) Get(id repository.LobbyID) (repository.ServerName, error) {
	serverName, ok := r.lobbies[id]
	if !ok {
		return "", domain.ErrLobbyNotFound
	}
	return serverName, nil
}
