// Package repository defines the interfaces for the repositories used in the application.
// It provides an abstraction layer for data access, allowing for different implementations without affecting the business logic.
package repository

import (
	"time"

	"github.com/sceredi/co-type/common/domain"
)

// ServerRepository defines the interface for managing servers in the broker.
type ServerRepository interface {
	// Create adds a new server to the repository. It returns the created server or an error if the operation fails.
	Create(server *domain.Server) (*domain.Server, error)
	// List retrieves all servers from the repository.
	List() []*domain.Server
	// UpdateLastSeen records the time a server was last seen alive and updates its load.
	UpdateLastSeen(name string, load int, t time.Time) error
	// Delete removes a server by name. Returns ErrServerNotFound if no such server exists.
	Delete(name string) error
	// ListStale returns the names of servers whose last heartbeat is older than ttl.
	ListStale(ttl time.Duration) []string
}

// LobbyID is a type alias for string, representing the unique identifier of a lobby.
type LobbyID string

// ServerName is a type alias for string, representing the name of a server managing a lobby.
type ServerName string

// LobbyRepository defines the interface for managing lobbies in the broker.
type LobbyRepository interface {
	// Create registers a new lobby in the repository with the server name managing it. It returns an error if the lobby already exists or if the operation fails.
	Create(lobbyID LobbyID, serverName ServerName) error
	// Get retrieves a lobby by its ID from the repository. It returns the server name managing the lobby or an error if the lobby is not found or if the operation fails.
	Get(id LobbyID) (ServerName, error)
	// Delete removes a lobby by its ID from the repository. It returns an error if the lobby is not found or if the operation fails.
	Delete(id LobbyID) error
}
