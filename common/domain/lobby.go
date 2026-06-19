package domain

import (
	"errors"

	"github.com/sceredi/co-type/common/proto/lobby"
)

// ErrLobbyAlreadyExists is returned when a lobby with the same ID already exists.
var ErrLobbyAlreadyExists = errors.New("lobby already exists")

// ErrLobbyNotFound is returned when a lobby with the given ID is not found.
var ErrLobbyNotFound = errors.New("lobby not found")

// ErrPlayerNotInLobby is returned when a player is not in the lobby.
var ErrPlayerNotInLobby = errors.New("player not in lobby")

// ErrPlayerAlreadyInLobby is returned when a player with the same name is already in the lobby.
var ErrPlayerAlreadyInLobby = errors.New("player already in lobby")

// Lobby represents the state of the lobby.
type Lobby struct {
	ID      string
	Players []*Player
	Snippet string
}

// NewLobby creates a new lobby with the given ID and first player.
func NewLobby(id string, firstPlayer *Player) *Lobby {
	return &Lobby{
		ID:      id,
		Players: []*Player{firstPlayer},
		Snippet: "",
	}
}

// NewLobbyFromGRPC converts a gRPC Lobby message to the Lobby struct.
func NewLobbyFromGRPC(l *lobby.Lobby) *Lobby {
	players := make([]*Player, len(l.Players))
	for i, p := range l.Players {
		players[i] = &Player{
			Name:              p.Name,
			IsReady:           p.IsReady,
			AllowedCharacters: p.AllowedCharacters,
			BlockedCharacters: p.BlockedCharacters,
			CanDelete:         p.CanDelete,
		}
	}
	return &Lobby{
		ID:      l.Id,
		Players: players,
		Snippet: l.Snippet,
	}
}

// AddPlayers adds the given players to the lobby.
func (l *Lobby) AddPlayers(players ...*Player) {
	l.Players = append(l.Players, players...)
}

// UpdatePlayer updates the settings of the player with the given name. Returns true if the player was found and updated.
func (l *Lobby) UpdatePlayer(name string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) bool {
	for _, p := range l.Players {
		if p.Name == name {
			p.IsReady = isReady
			p.AllowedCharacters = allowedCharacters
			p.BlockedCharacters = blockedCharacters
			p.CanDelete = canDelete
			return true
		}
	}
	return false
}

// RemovePlayer removes the player with the given name from the lobby. Returns true if the player was found and removed.
func (l *Lobby) RemovePlayer(name string) bool {
	for i, p := range l.Players {
		if p.Name == name {
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			return true
		}
	}
	return false
}

// ToGRPCLobby converts the Lobby struct to its gRPC representation.
func (l *Lobby) ToGRPCLobby() *lobby.Lobby {
	players := make([]*lobby.Player, len(l.Players))
	for i, p := range l.Players {
		players[i] = &lobby.Player{
			Name:              p.Name,
			IsReady:           p.IsReady,
			AllowedCharacters: p.AllowedCharacters,
			BlockedCharacters: p.BlockedCharacters,
			CanDelete:         p.CanDelete,
		}
	}
	return &lobby.Lobby{
		Id:      l.ID,
		Players: players,
		Snippet: l.Snippet,
	}
}
