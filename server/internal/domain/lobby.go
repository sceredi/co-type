// Package domain contains the domain of the server service.
package domain

import "github.com/sceredi/co-type/common/domain"

// Lobby represents a game lobby in the server service.
type Lobby struct {
	Subs map[string]chan *Lobby
	Base *domain.Lobby
	// DisconnectedPlayers tracks players who dropped mid-game and are expected to reconnect.
	DisconnectedPlayers map[string]bool
}

// NewLobby creates a new instance of Lobby with the provided ID and first player, and returns it.
func NewLobby(id string, firstPlayer *domain.Player) *Lobby {
	base := domain.NewLobby(id, firstPlayer)
	return &Lobby{
		Base:                base,
		Subs:                make(map[string]chan *Lobby),
		DisconnectedPlayers: make(map[string]bool),
	}
}
