package domain

import "time"

// GameStats represents the statistics of a game.
type GameStats struct {
	TotalTime time.Duration
	Lobby     Lobby
}
