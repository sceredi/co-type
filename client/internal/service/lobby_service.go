package service

import (
	"errors"

	"github.com/sceredi/co-type/client/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/common/domain"
)

// ErrLobbyGatewayNotSet is returned when the LobbyService is used without setting a gateway first.
var ErrLobbyGatewayNotSet = errors.New("lobby gateway not set")

// LobbyService defines the interface for managing lobbies.
type LobbyService interface {
	// Create creates a new lobby with the given id and host name. It returns the created lobby or an error if something goes wrong.
	Create(id, hostName string) (*domain.Lobby, error)

	// Join tries to join the lobby with the given id and player name. It returns the joined lobby or an error if something goes wrong.
	Join(id, playerName string) (*domain.Lobby, error)

	// Leave leaves the lobby with the given id for the given player name. It returns an error if something goes wrong.
	Leave(id, playerName string) error

	// EditPlayer edits the settings of the player with the given name in the lobby. It returns the updated lobby or an error if something goes wrong.
	EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error)

	// Connect connects to the given server. It returns an error if the connection fails.
	Connect(target *domain.Server) error

	// Subscribe subscribes to lobby events and returns a channel that receives updated lobby state.
	Subscribe(lobbyID, playerName string) (<-chan *domain.Lobby, error)
}

type lobbyService struct {
	gtw gateway.LobbyGateway
}

// NewLobbyService creates a new instance of LobbyService.
func NewLobbyService(gtw gateway.LobbyGateway) LobbyService {
	return &lobbyService{gtw: gtw}
}

func (s *lobbyService) Create(id, hostName string) (*domain.Lobby, error) {
	if s.gtw == nil {
		return nil, ErrLobbyGatewayNotSet
	}
	return s.gtw.Create(id, hostName)
}

func (s *lobbyService) Join(id, playerName string) (*domain.Lobby, error) {
	if s.gtw == nil {
		return nil, ErrLobbyGatewayNotSet
	}
	return s.gtw.Join(id, playerName)
}

func (s *lobbyService) Leave(id, playerName string) error {
	if s.gtw == nil {
		return ErrLobbyGatewayNotSet
	}
	return s.gtw.Leave(id, playerName)
}

func (s *lobbyService) EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error) {
	if s.gtw == nil {
		return nil, ErrLobbyGatewayNotSet
	}
	return s.gtw.EditPlayer(lobbyID, playerName, isReady, allowedCharacters, blockedCharacters, canDelete)
}

func (s *lobbyService) Connect(target *domain.Server) error {
	return s.gtw.Connect(target)
}

func (s *lobbyService) Subscribe(lobbyID, playerName string) (<-chan *domain.Lobby, error) {
	if s.gtw == nil {
		return nil, ErrLobbyGatewayNotSet
	}
	return s.gtw.Subscribe(lobbyID, playerName)
}
