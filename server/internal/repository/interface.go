// Package repository contains the repository layer of the server service.
package repository

import "github.com/sceredi/co-type/common/domain"

// LobbyRepository defines the interface for managing lobbies.
type LobbyRepository interface {
	// Create saves a new lobby to the repository and returns the created lobby or an error if the operation fails.
	Create(lobby domain.Lobby) (*domain.Lobby, error)
}
