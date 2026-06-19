package service

import (
	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

// LobbyService defines the interface for managing lobbies in the broker service.
type LobbyService interface {
	// Create registers a new lobby in the broker service. It takes a lobby ID and the server name managing it, and returns an error if the operation fails (e.g., if the lobby already exists).
	Create(lobbyID repository.LobbyID, serverName repository.ServerName) error
	// Get retrieves a lobby by its ID from the broker service. It returns the server managing the lobby or an error if the lobby is not found or if the operation fails.
	Get(lobbyID repository.LobbyID) (*domain.Server, error)
}

type lobbyService struct {
	lobbyRepo repository.LobbyRepository
	serverSvc ServerService
}

// NewLobbyService creates a new instance of LobbyService with the provided LobbyRepository and ServerService, and returns it.
func NewLobbyService(lobbyRepo repository.LobbyRepository, serverSvc ServerService) LobbyService {
	return &lobbyService{
		lobbyRepo: lobbyRepo,
		serverSvc: serverSvc,
	}
}

func (s *lobbyService) Create(lobbyID repository.LobbyID, serverName repository.ServerName) error {
	return s.lobbyRepo.Create(lobbyID, serverName)
}

func (s *lobbyService) Get(lobbyID repository.LobbyID) (*domain.Server, error) {
	serverName, err := s.lobbyRepo.Get(lobbyID)
	if err != nil {
		return nil, err
	}
	server, err := s.serverSvc.GetByName(string(serverName))
	if err != nil {
		return nil, err
	}
	return server, nil
}
