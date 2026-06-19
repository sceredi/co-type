package service

import (
	"context"
	"log/slog"

	commondomain "github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/server/internal/domain"
	"github.com/sceredi/co-type/server/internal/repository"
)

// LobbyService defines the interface for managing lobbies in the server service.
type LobbyService interface {
	Create(id, userName string) (*domain.Lobby, error)
	Join(id, userName string) (*domain.Lobby, error)
	Leave(id, userName string) error
	EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error)
	Get(id string) *domain.Lobby
}

type lobbyService struct {
	serverName string
	controlGtw gateway.ControlGateway
	lobbyRepo  repository.LobbyRepository
}

// NewLobbyService creates a new instance of LobbyService with the provided ControlGateway and LobbyRepository and returns it.
func NewLobbyService(serverName string, controlGtw gateway.ControlGateway, lobbyRepo repository.LobbyRepository) LobbyService {
	return &lobbyService{serverName: serverName, controlGtw: controlGtw, lobbyRepo: lobbyRepo}
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

	err = s.controlGtw.RegisterLobby(id, s.serverName)
	if err != nil {
		if deleteErr := s.lobbyRepo.Delete(id); deleteErr != nil {
			slog.ErrorContext(context.Background(), "Failed to delete lobby after control gateway registration failure",
				slog.String("id", id),
				slog.String("error", deleteErr.Error()),
			)
		}
		return nil, err
	}

	ch := make(chan *domain.Lobby, 64)
	lobby.Subs[userName] = ch
	return lobby, nil
}

func (s *lobbyService) Join(id, userName string) (*domain.Lobby, error) {
	slog.DebugContext(context.Background(), "Joining lobby",
		slog.String("id", id),
		slog.String("userName", userName),
	)
	lobby := s.lobbyRepo.Get(id)
	if lobby == nil {
		return nil, commondomain.ErrLobbyNotFound
	}
	player := commondomain.NewPlayer(userName)
	for _, p := range lobby.Base.Players {
		if p.Name == userName {
			return nil, commondomain.ErrPlayerAlreadyInLobby
		}
	}
	lobby.Base.AddPlayers(player)

	for _, ch := range lobby.Subs {
		ch <- lobby
	}

	ch := make(chan *domain.Lobby, 64)
	lobby.Subs[userName] = ch
	return lobby, nil
}

func (s *lobbyService) Get(id string) *domain.Lobby {
	return s.lobbyRepo.Get(id)
}

func (s *lobbyService) EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error) {
	slog.DebugContext(context.Background(), "Editing player",
		slog.String("lobbyID", lobbyID),
		slog.String("playerName", playerName),
	)
	lobby := s.lobbyRepo.Get(lobbyID)
	if lobby == nil {
		return nil, commondomain.ErrLobbyNotFound
	}
	updated := lobby.Base.UpdatePlayer(playerName, isReady, allowedCharacters, blockedCharacters, canDelete)
	if !updated {
		return nil, commondomain.ErrPlayerNotInLobby
	}

	for _, ch := range lobby.Subs {
		ch <- lobby
	}
	return lobby, nil
}

func (s *lobbyService) Leave(id, userName string) error {
	slog.DebugContext(context.Background(), "Leaving lobby",
		slog.String("id", id),
		slog.String("userName", userName),
	)
	lobby := s.lobbyRepo.Get(id)
	if lobby == nil {
		return commondomain.ErrLobbyNotFound
	}
	removed := lobby.Base.RemovePlayer(userName)
	if !removed {
		return commondomain.ErrPlayerNotInLobby
	}

	if ch, ok := lobby.Subs[userName]; ok {
		close(ch)
		delete(lobby.Subs, userName)
	}

	if len(lobby.Base.Players) == 0 {
		if err := s.lobbyRepo.Delete(id); err != nil {
			slog.ErrorContext(context.Background(), "Failed to delete empty lobby from repository",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		if err := s.controlGtw.UnregisterLobby(id); err != nil {
			slog.ErrorContext(context.Background(), "Failed to unregister empty lobby from broker",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		return nil
	}

	for _, ch := range lobby.Subs {
		ch <- lobby
	}
	return nil
}
