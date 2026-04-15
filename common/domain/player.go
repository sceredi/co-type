package domain

// Player represents a player in the lobby.
type Player struct {
	Name              string
	IsReady           bool
	AllowedCharacters string
	BlockedCharacters string
	CanDelete         bool
}

func NewPlayer(name string) Player {
	return Player{Name: name, IsReady: false}
}
