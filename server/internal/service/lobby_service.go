package service

import (
	"context"
	"log/slog"

	"github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/repository"
)

// LobbyService defines the interface for managing lobbies in the server service.
type LobbyService interface {
	Create(id, userName string) (*domain.Lobby, error)
}

type lobbyService struct {
	lobbyRepo repository.LobbyRepository
}

// NewLobbyService creates a new instance of LobbyService with the provided LobbyRepository and returns it.
func NewLobbyService(lobbyRepo repository.LobbyRepository) LobbyService {
	return &lobbyService{lobbyRepo: lobbyRepo}
}

func (s *lobbyService) Create(id, userName string) (*domain.Lobby, error) {
	slog.DebugContext(context.Background(), "Creating lobby",
		slog.String("id", id),
		slog.String("userName", userName),
	)
	host := domain.NewPlayer(userName)
	lobby := domain.NewLobby(id, host)
	return s.lobbyRepo.Create(lobby)
}
