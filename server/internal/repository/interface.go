// Package repository contains the repository layer of the server service.
package repository

import "github.com/sceredi/co-type/server/internal/domain"

// LobbyRepository defines the interface for managing lobbies.
type LobbyRepository interface {
	// Create saves a new lobby to the repository and returns the created lobby or an error if the operation fails.
	Create(lobby *domain.Lobby) (*domain.Lobby, error)
	// Get retrieves a lobby by its ID from the repository. It returns the lobby if found, or nil if no lobby with the given ID exists.
	Get(id string) *domain.Lobby
	// Delete removes a lobby by its ID from the repository. It returns an error if the operation fails.
	Delete(id string) error
	// ListIDs returns all lobby IDs currently stored in the repository.
	ListIDs() []string
}
