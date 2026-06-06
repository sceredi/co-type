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

// Lobby represents the state of the lobby.
type Lobby struct {
	ID      string
	Players []*Player
	Host    *Player
	Snippet string
}

// NewLobby creates a new lobby with the given ID and host player.
func NewLobby(id string, host *Player) *Lobby {
	return &Lobby{
		ID:      id,
		Players: []*Player{host},
		Host:    host,
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
		Host: &Player{
			Name:              l.Host.Name,
			IsReady:           l.Host.IsReady,
			AllowedCharacters: l.Host.AllowedCharacters,
			BlockedCharacters: l.Host.BlockedCharacters,
			CanDelete:         l.Host.CanDelete,
		},
		Snippet: l.Snippet,
	}
}

// AddPlayers adds the given players to the lobby.
func (l *Lobby) AddPlayers(players ...*Player) {
	l.Players = append(l.Players, players...)
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
		Host: &lobby.Player{
			Name:              l.Host.Name,
			IsReady:           l.Host.IsReady,
			AllowedCharacters: l.Host.AllowedCharacters,
			BlockedCharacters: l.Host.BlockedCharacters,
			CanDelete:         l.Host.CanDelete,
		},
		Snippet: l.Snippet,
	}
}
