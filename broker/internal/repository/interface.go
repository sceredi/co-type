// Package repository defines the interfaces for the repositories used in the application.
// It provides an abstraction layer for data access, allowing for different implementations without affecting the business logic.
package repository

import "github.com/sceredi/co-type/common/domain"

// ServerRepository defines the interface for managing servers in the broker.
type ServerRepository interface {
	// Create adds a new server to the repository. It returns the created server or an error if the operation fails.
	Create(server *domain.Server) (*domain.Server, error)
	// List retrieves all servers from the repository.
	List() []*domain.Server
}
