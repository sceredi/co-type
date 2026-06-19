// Package domain contains the domain of the server service.
package domain

import "github.com/sceredi/co-type/common/domain"

// Lobby represents a game lobby in the server service.
type Lobby struct {
	Subs map[string]chan *Lobby
	Base *domain.Lobby
}

// NewLobby creates a new instance of Lobby with the provided ID and first player, and returns it.
func NewLobby(id string, firstPlayer *domain.Player) *Lobby {
	base := domain.NewLobby(id, firstPlayer)
	subs := make(map[string]chan *Lobby)
	return &Lobby{Base: base, Subs: subs}
}
