package service

import (
	"context"
	"log/slog"

	commondomain "github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/server/internal/domain"
	"github.com/sceredi/co-type/server/internal/repository"
	"github.com/sceredi/co-type/server/internal/snippet"
)

// LobbyService defines the interface for managing lobbies in the server service.
type LobbyService interface {
	Create(id, userName string) (*domain.Lobby, error)
	Join(id, userName string) (*domain.Lobby, error)
	Leave(id, userName string) error
	// PlayerDisconnected marks a player as disconnected mid-game.
	PlayerDisconnected(id, userName string) error
	EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error)
	Ready(lobbyID, playerName string) (*domain.Lobby, error)
	Get(id string) *domain.Lobby
	SendKeyPress(lobbyID, playerName, key string, isBackspace bool) (*domain.Lobby, error)
	// ListIDs returns the IDs of all lobbies currently managed by this server.
	ListIDs() []string
	// ResumeGame restores a crashed game on this server.
	ResumeGame(lobby *commondomain.Lobby, playerName string) (*domain.Lobby, error)
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

func reconnectPlayer(id, userName string, lobby *domain.Lobby) (*domain.Lobby, error) {
	slog.DebugContext(context.Background(), "Player reconnecting to paused game",
		slog.String("id", id),
		slog.String("userName", userName),
	)
	delete(lobby.DisconnectedPlayers, userName)
	ch := make(chan *domain.Lobby, 64)
	lobby.Subs[userName] = ch
	if len(lobby.DisconnectedPlayers) == 0 {
		lobby.Base.Status = commondomain.LobbyPlaying
	}
	for _, ch := range lobby.Subs {
		ch <- lobby
	}
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

	if lobby.DisconnectedPlayers[userName] {
		return reconnectPlayer(id, userName, lobby)
	}

	if lobby.Base.Status != commondomain.LobbyWaitingForPlayers {
		return nil, commondomain.ErrLobbyNotOpen
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

	// Pause the game if a player leaves mid-game.
	if lobby.Base.Status == commondomain.LobbyPlaying {
		lobby.Base.Status = commondomain.LobbyPaused
	}

	for _, ch := range lobby.Subs {
		ch <- lobby
	}
	return nil
}

func (s *lobbyService) Ready(lobbyID, playerName string) (*domain.Lobby, error) {
	slog.DebugContext(context.Background(), "Toggling player ready",
		slog.String("lobbyID", lobbyID),
		slog.String("playerName", playerName),
	)
	lobby := s.lobbyRepo.Get(lobbyID)
	if lobby == nil {
		return nil, commondomain.ErrLobbyNotFound
	}

	var player *commondomain.Player
	for _, p := range lobby.Base.Players {
		if p.Name == playerName {
			player = p
			break
		}
	}
	if player == nil {
		return nil, commondomain.ErrPlayerNotInLobby
	}

	player.IsReady = !player.IsReady

	if allPlayersReady(lobby.Base) {
		lobby.Base.Game = *commondomain.NewGameInfo(snippet.Random())
		lobby.Base.Status = commondomain.LobbyPlaying
	}

	for _, ch := range lobby.Subs {
		ch <- lobby
	}
	return lobby, nil
}

func allPlayersReady(l *commondomain.Lobby) bool {
	if len(l.Players) == 0 {
		return false
	}
	for _, p := range l.Players {
		if !p.IsReady {
			return false
		}
	}
	return true
}

func (s *lobbyService) SendKeyPress(lobbyID, playerName, key string, isBackspace bool) (*domain.Lobby, error) {
	slog.DebugContext(context.Background(), "Sending key press",
		slog.String("lobbyID", lobbyID),
		slog.String("playerName", playerName),
		slog.String("key", key),
		slog.Bool("isBackspace", isBackspace),
	)
	lobby := s.lobbyRepo.Get(lobbyID)
	if lobby == nil {
		return nil, commondomain.ErrLobbyNotFound
	}
	if lobby.Base.Status != commondomain.LobbyPlaying {
		return nil, commondomain.ErrGameNotPlaying
	}

	var player *commondomain.Player
	for _, p := range lobby.Base.Players {
		if p.Name == playerName {
			player = p
			break
		}
	}
	if player == nil {
		return nil, commondomain.ErrPlayerNotInLobby
	}

	if err := validateKey(player, key, isBackspace); err != nil {
		return nil, err
	}

	ended := applyKeyPress(&lobby.Base.Game, key, isBackspace)
	if ended {
		lobby.Base.Status = commondomain.LobbyGameEnded
	}

	for _, ch := range lobby.Subs {
		ch <- lobby
	}
	return lobby, nil
}

func (s *lobbyService) PlayerDisconnected(id, userName string) error {
	slog.DebugContext(context.Background(), "Player disconnected mid-game",
		slog.String("id", id),
		slog.String("userName", userName),
	)
	lobby := s.lobbyRepo.Get(id)
	if lobby == nil {
		return commondomain.ErrLobbyNotFound
	}
	lobby.DisconnectedPlayers[userName] = true
	if ch, ok := lobby.Subs[userName]; ok {
		close(ch)
		delete(lobby.Subs, userName)
	}

	if len(lobby.Subs) == 0 {
		if err := s.lobbyRepo.Delete(id); err != nil {
			slog.ErrorContext(context.Background(), "Failed to delete abandoned lobby from repository",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		if err := s.controlGtw.UnregisterLobby(id); err != nil {
			slog.ErrorContext(context.Background(), "Failed to unregister abandoned lobby from broker",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		return nil
	}

	if lobby.Base.Status == commondomain.LobbyPlaying {
		lobby.Base.Status = commondomain.LobbyPaused
	}
	for _, ch := range lobby.Subs {
		ch <- lobby
	}
	return nil
}

func (s *lobbyService) ListIDs() []string {
	return s.lobbyRepo.ListIDs()
}

func (s *lobbyService) ResumeGame(incomingLobby *commondomain.Lobby, playerName string) (*domain.Lobby, error) {
	slog.InfoContext(context.Background(), "Resuming game",
		slog.String("id", incomingLobby.ID),
		slog.String("playerName", playerName),
	)

	existing := s.lobbyRepo.Get(incomingLobby.ID)
	if existing == nil {
		return s.initResumedLobby(incomingLobby, playerName)
	}
	return s.reconnectPlayer(existing, incomingLobby, playerName)
}

func (s *lobbyService) initResumedLobby(incomingLobby *commondomain.Lobby, playerName string) (*domain.Lobby, error) {
	lobby := &domain.Lobby{
		Base:                incomingLobby,
		Subs:                make(map[string]chan *domain.Lobby),
		DisconnectedPlayers: make(map[string]bool),
	}
	lobby.Base.Status = commondomain.LobbyPaused
	s.markOthersDisconnected(lobby, playerName)

	if _, err := s.lobbyRepo.Create(lobby); err != nil {
		return nil, err
	}
	if err := s.controlGtw.RegisterLobby(incomingLobby.ID, s.serverName); err != nil {
		_ = s.lobbyRepo.Delete(incomingLobby.ID)
		return nil, err
	}

	slog.InfoContext(context.Background(), "Lobby created for resume, waiting for other players",
		slog.String("id", incomingLobby.ID),
		slog.Int("disconnected", len(lobby.DisconnectedPlayers)),
	)
	ch := make(chan *domain.Lobby, 64)
	lobby.Subs[playerName] = ch
	return lobby, nil
}

func (s *lobbyService) markOthersDisconnected(lobby *domain.Lobby, playerName string) {
	for _, p := range lobby.Base.Players {
		if p.Name != playerName {
			lobby.DisconnectedPlayers[p.Name] = true
		}
	}
}

func (s *lobbyService) reconnectPlayer(existing *domain.Lobby, incomingLobby *commondomain.Lobby, playerName string) (*domain.Lobby, error) {
	s.mergeGameState(existing, incomingLobby)

	delete(existing.DisconnectedPlayers, playerName)
	ch := make(chan *domain.Lobby, 64)
	existing.Subs[playerName] = ch

	if len(existing.DisconnectedPlayers) == 0 {
		slog.InfoContext(context.Background(), "All players reconnected, resuming game", slog.String("id", incomingLobby.ID))
		existing.Base.Status = commondomain.LobbyPlaying
	} else {
		slog.InfoContext(context.Background(), "Player reconnected, still waiting",
			slog.String("id", incomingLobby.ID),
			slog.String("player", playerName),
			slog.Int("remaining_disconnected", len(existing.DisconnectedPlayers)),
		)
	}

	for _, sub := range existing.Subs {
		sub <- existing
	}
	return existing, nil
}

func (s *lobbyService) mergeGameState(existing *domain.Lobby, incoming *commondomain.Lobby) {
	if incoming.Game.Revision > existing.Base.Game.Revision {
		slog.InfoContext(context.Background(), "Updating game state to higher revision",
			slog.String("id", incoming.ID),
			slog.Int64("new_revision", incoming.Game.Revision),
		)
		existing.Base.Game = incoming.Game
	}
}
