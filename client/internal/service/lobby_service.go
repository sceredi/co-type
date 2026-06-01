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
	// SetGateway sets the gateway for the LobbyService. This method must be called before using any other method of the service.
	SetGateway(gtw gateway.LobbyGateway)
	// Create creates a new lobby with the given id and host name. It returns the created lobby or an error if something goes wrong.
	Create(id, hostName string) (*domain.Lobby, error)
}

type lobbyService struct {
	gtw gateway.LobbyGateway
}

// NewLobbyService creates a new instance of LobbyService.
func NewLobbyService() LobbyService {
	return &lobbyService{}
}

func (s *lobbyService) SetGateway(gtw gateway.LobbyGateway) {
	s.gtw = gtw
}

func (s *lobbyService) Create(id, hostName string) (*domain.Lobby, error) {
	if s.gtw == nil {
		return nil, ErrLobbyGatewayNotSet
	}
	return s.gtw.Create(id, hostName)
}
