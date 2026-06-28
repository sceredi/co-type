package service

import (
	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/common/domain"
)

// LobbyService defines the interface for managing lobbies in the broker service.
type LobbyService interface {
	// Create registers a new lobby in the broker service. It takes a lobby ID and the server name managing it, and returns an error if the operation fails (e.g., if the lobby already exists).
	Create(lobbyID repository.LobbyID, serverName repository.ServerName) error
	// Upsert registers or updates a lobby in the broker service. It never returns ErrLobbyAlreadyExists.
	Upsert(lobbyID repository.LobbyID, serverName repository.ServerName)
	// Get retrieves a lobby by its ID from the broker service. It returns the server managing the lobby or an error if the lobby is not found or if the operation fails.
	Get(lobbyID repository.LobbyID) (*domain.Server, error)
	// Delete removes a lobby by its ID from the broker service. It returns an error if the lobby is not found or if the operation fails.
	Delete(lobbyID repository.LobbyID) error
	// RequestResumeGame assigns a server to host the resumed lobby. If a server has already been
	// assigned for this lobby (a previous player called first), the same server is returned.
	// The assignment is stored as pending until the server calls RegisterLobby.
	RequestResumeGame(lobbyID repository.LobbyID) (*domain.Server, error)
	// RemoveFromPending removes a lobby from the pending-resume list. Called when the server
	// registers the lobby via RegisterLobby.
	RemoveFromPending(lobbyID repository.LobbyID)
}

type lobbyService struct {
	lobbyRepo         repository.LobbyRepository
	pendingResumeRepo repository.PendingResumeRepository
	serverSvc         ServerService
}

// NewLobbyService creates a new instance of LobbyService with the provided LobbyRepository and ServerService, and returns it.
func NewLobbyService(lobbyRepo repository.LobbyRepository, pendingResumeRepo repository.PendingResumeRepository, serverSvc ServerService) LobbyService {
	return &lobbyService{
		lobbyRepo:         lobbyRepo,
		pendingResumeRepo: pendingResumeRepo,
		serverSvc:         serverSvc,
	}
}

func (s *lobbyService) Create(lobbyID repository.LobbyID, serverName repository.ServerName) error {
	return s.lobbyRepo.Create(lobbyID, serverName)
}

func (s *lobbyService) Upsert(lobbyID repository.LobbyID, serverName repository.ServerName) {
	s.lobbyRepo.Upsert(lobbyID, serverName)
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

func (s *lobbyService) Delete(lobbyID repository.LobbyID) error {
	return s.lobbyRepo.Delete(lobbyID)
}

func (s *lobbyService) RequestResumeGame(lobbyID repository.LobbyID) (*domain.Server, error) {
	if assignedName, ok := s.pendingResumeRepo.Get(lobbyID); ok {
		return s.serverSvc.GetByName(string(assignedName))
	}

	var crashedServer string
	if prevServer, err := s.lobbyRepo.Get(lobbyID); err == nil {
		crashedServer = string(prevServer)
	}

	server, err := s.serverSvc.LowestLoadExcluding(crashedServer)
	if err != nil {
		return nil, err
	}
	s.pendingResumeRepo.Set(lobbyID, repository.ServerName(server.Name))
	return server, nil
}

func (s *lobbyService) RemoveFromPending(lobbyID repository.LobbyID) {
	s.pendingResumeRepo.Delete(lobbyID)
}
