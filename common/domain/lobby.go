package domain

import "errors"

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

// AddPlayers adds the given players to the lobby.
func (l *Lobby) AddPlayers(players ...*Player) {
	l.Players = append(l.Players, players...)
}

// LobbyEvent represents an event that can occur in the lobby.
type LobbyEvent interface {
	isLobbyEvent()
}

// PlayerJoin represents an event where a player joins the lobby.
type PlayerJoin struct {
	PlayerName string
}

// PlayerLeave represents an event where a player leaves the lobby.
type PlayerLeave struct {
	PlayerName string
}

// PlayerReady represents an event where a player changes their ready status in the lobby.
type PlayerReady struct {
	PlayerName string
	IsReady    bool
}

// PlayerEdit represents an event where a player edits their allowed/blocked characters or delete permission in the lobby.
type PlayerEdit struct {
	PlayerName        string
	AllowedCharacters string
	BlockedCharacters string
	CanDelete         bool
}

func (PlayerJoin) isLobbyEvent()  {}
func (PlayerLeave) isLobbyEvent() {}
func (PlayerReady) isLobbyEvent() {}
func (PlayerEdit) isLobbyEvent()  {}
