// Package service contains the service layer of the server service.
// The service layer is responsible for implementing the business logic of the server service.
// It defines the interfaces that the handlers will use to interact with the domain models and repositories.
package service

import "github.com/sceredi/co-type/server/internal/api/grpc/gateway"

// DiscoveryService defines the interface for managing service discovery in the server service.
// It provides a method to register the server with the service discovery mechanism.
// The Register method is responsible for registering the server and returns an error if the operation fails.
type DiscoveryService interface {
	// Register registers the server with the service discovery mechanism using the provided host and port.
	Register(name, host string, port int) error
}

type discoveryService struct {
	gtw gateway.ControlGateway
}

// NewDiscoveryService creates a new instance of DiscoveryService with the provided ControlGateway and returns it.
func NewDiscoveryService(gtw gateway.ControlGateway) DiscoveryService {
	return &discoveryService{gtw: gtw}
}

func (s *discoveryService) Register(name, host string, port int) error {
	return s.gtw.Register(name, host, port)
}
