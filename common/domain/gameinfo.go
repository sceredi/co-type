package domain

// GameInfo holds the shared game state that is propagated to all clients via the lobby subscription.
// Revision is a monotonically increasing counter incremented on every mutation; clients can use it
// to reconstruct the most recent state after a disconnection.
type GameInfo struct {
	Snippet      string
	CorrectChars int64
	WrongChars   int64
	Revision     int64
}

// NewGameInfo creates a new GameInfo for a game starting with the given snippet.
// CorrectChars, WrongChars and Revision are initialised to zero.
func NewGameInfo(snippet string) *GameInfo {
	return &GameInfo{Snippet: snippet}
}
