// Package service contains the service layer of the client.
package service

import (
	"github.com/sceredi/co-type/client/internal/grpc"
	"github.com/sceredi/co-type/common/domain"
)

// DiscoveryService defines the interface for managing server discovery.
type DiscoveryService interface {
	GetAvailableServer() (*domain.Server, error)
}

type discoveryService struct {
	gtw grpc.DiscoveryGateway
}

// NewDiscoveryService creates a new instance of DiscoveryService with the provided DiscoveryGateway and returns it.
func NewDiscoveryService(gtw grpc.DiscoveryGateway) DiscoveryService {
	return &discoveryService{gtw: gtw}
}

func (s *discoveryService) GetAvailableServer() (*domain.Server, error) {
	return s.gtw.AvailableServer()
}
