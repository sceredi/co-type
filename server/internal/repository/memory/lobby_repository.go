// Package memory provides an in-memory implementation of the repository interfaces.
package memory

import (
	commondomain "github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/domain"
	"github.com/sceredi/co-type/server/internal/repository"
)

type lobbyRepository struct {
	lobbies map[string]*domain.Lobby
}

// NewLobbyRepository creates a new instance of an in-memory LobbyRepository.
func NewLobbyRepository() repository.LobbyRepository {
	return &lobbyRepository{lobbies: make(map[string]*domain.Lobby)}
}

func (r *lobbyRepository) Create(lobby *domain.Lobby) (*domain.Lobby, error) {
	_, ok := r.lobbies[lobby.Base.ID]
	if ok {
		return nil, commondomain.ErrLobbyAlreadyExists
	}
	r.lobbies[lobby.Base.ID] = lobby
	return lobby, nil
}

func (r *lobbyRepository) Get(id string) *domain.Lobby {
	lobby, ok := r.lobbies[id]
	if !ok {
		return nil
	}
	return lobby
}
