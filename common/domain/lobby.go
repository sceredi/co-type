package domain

const MaxPlayers = 8

// Lobby represents the state of the lobby.
type Lobby struct {
	ID         string
	Players    []Player
	Host       Player
	MaxPlayers int
}

func NewLobby(id string, host Player) Lobby {
	return Lobby{
		ID:         id,
		Players:    []Player{host},
		Host:       host,
		MaxPlayers: MaxPlayers,
	}
}
