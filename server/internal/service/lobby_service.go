package service

import (
	"context"
	"log/slog"

	commondomain "github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/domain"
	"github.com/sceredi/co-type/server/internal/repository"
)

// LobbyService defines the interface for managing lobbies in the server service.
type LobbyService interface {
	Create(id, userName string) (*domain.Lobby, error)
	Get(id string) *domain.Lobby
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
	host := commondomain.NewPlayer(userName)
	lobby := domain.NewLobby(id, host)
	lobby, err := s.lobbyRepo.Create(lobby)
	if err != nil {
		return nil, err
	}
	ch := make(chan commondomain.LobbyEvent, 64)
	lobby.Subs[userName] = ch
	return lobby, nil
}

func (s *lobbyService) Get(id string) *domain.Lobby {
	return s.lobbyRepo.Get(id)
}
