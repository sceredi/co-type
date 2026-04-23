// Package service contains the service layer of the broker service.
// The service layer is responsible for implementing the business logic of the broker service.
// It defines the interfaces that the handlers will use to interact with the domain models and repositories.
package service

import (
	"github.com/sceredi/co-type/broker/internal/domain"
	"github.com/sceredi/co-type/broker/internal/repository"
)

// ServerService defines the interface for managing servers in the broker service.
type ServerService interface {
	// Create creates a new server in the broker service. It takes a CreateServerRequest and returns the created Server or an error if the operation fails.
	Create(req domain.CreateServerRequest) (*domain.Server, error)
}

type serverService struct {
	serverRepo repository.ServerRepository
}

// NewServerService creates a new instance of ServerService with the provided ServerRepository and returns it.
func NewServerService(serverRepo repository.ServerRepository) ServerService {
	return &serverService{serverRepo: serverRepo}
}

func (s *serverService) Create(req domain.CreateServerRequest) (*domain.Server, error) {
	server, err := domain.NewServer(req.Addr)
	if err != nil {
		return nil, err
	}
	result, err := s.serverRepo.Create(server)
	if err != nil {
		return nil, err
	}
	return result, nil
}
