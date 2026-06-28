package memory

import (
	"github.com/sceredi/co-type/broker/internal/repository"
)

type pendingResumeRepository struct {
	pending map[repository.LobbyID]repository.ServerName
}

// NewPendingResumeRepository creates a new in-memory PendingResumeRepository.
func NewPendingResumeRepository() repository.PendingResumeRepository {
	return &pendingResumeRepository{
		pending: make(map[repository.LobbyID]repository.ServerName),
	}
}

func (r *pendingResumeRepository) Set(lobbyID repository.LobbyID, serverName repository.ServerName) {
	r.pending[lobbyID] = serverName
}

func (r *pendingResumeRepository) Get(lobbyID repository.LobbyID) (repository.ServerName, bool) {
	name, ok := r.pending[lobbyID]
	return name, ok
}

func (r *pendingResumeRepository) Delete(lobbyID repository.LobbyID) {
	delete(r.pending, lobbyID)
}
