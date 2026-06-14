// Package service contains the service layer of the client.
package service

import (
	"github.com/sceredi/co-type/client/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/common/domain"
)

// DiscoveryService defines the interface for managing server discovery.
type DiscoveryService interface {
	GetAvailableServer() (*domain.Server, error)
	GetHostByLobby(lobbyCode string) (*domain.Server, error)
}

type discoveryService struct {
	gtw gateway.DiscoveryGateway
}

// NewDiscoveryService creates a new instance of DiscoveryService with the provided DiscoveryGateway and returns it.
func NewDiscoveryService(gtw gateway.DiscoveryGateway) DiscoveryService {
	return &discoveryService{gtw: gtw}
}

func (s *discoveryService) GetAvailableServer() (*domain.Server, error) {
	return s.gtw.AvailableServer()
}

func (s *discoveryService) GetHostByLobby(lobbyCode string) (*domain.Server, error) {
	return s.gtw.HostByLobby(lobbyCode)
}
