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

func (r *lobbyRepository) Delete(id string) error {
	_, ok := r.lobbies[id]
	if !ok {
		return commondomain.ErrLobbyNotFound
	}
	delete(r.lobbies, id)
	return nil
}

func (r *lobbyRepository) ListIDs() []string {
	ids := make([]string, 0, len(r.lobbies))
	for id := range r.lobbies {
		ids = append(ids, id)
	}
	return ids
}
