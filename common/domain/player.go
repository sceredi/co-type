// Package domain contains the domain models for the co-type game.
package domain

// Player represents a player in the lobby.
type Player struct {
	Name              string
	IsReady           bool
	AllowedCharacters string
	BlockedCharacters string
	CanDelete         bool
}

// NewPlayer creates a new player with the given name and default values for other fields.
func NewPlayer(name string) Player {
	return Player{Name: name, IsReady: false}
}
