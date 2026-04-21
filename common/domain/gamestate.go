package domain

// Game represents the ongoing game, including the lobby and the current state of the game.
type Game struct {
	Lobby Lobby
	State GameState
}

// GameStatus represents the current status of the game.
type GameStatus int

const (
	// Paused indicates that the game is currently paused.
	Paused GameStatus = iota
	// Running indicates that the game is currently running.
	Running
	// Terminated indicates that the game has been terminated.
	Terminated
)

// GameState represents the current state of the game.
type GameState struct {
	CorrectChars int
	LastChar     int
	Status       GameStatus
}
