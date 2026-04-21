package domain

// Lobby represents the state of the lobby.
type Lobby struct {
	ID      string
	Players []*Player
	Host    *Player
	Snippet string
}

// NewLobby creates a new lobby with the given ID and host player.
func NewLobby(id string, host *Player) Lobby {
	return Lobby{
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
