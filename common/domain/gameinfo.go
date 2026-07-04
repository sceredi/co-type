package domain

// GameInfo holds the shared game state that is propagated to all clients via the lobby subscription.
// Revision is a monotonically increasing counter incremented on every mutation; clients can use it
// to reconstruct the most recent state after a disconnection.
// StartTimeMs is the Unix timestamp in milliseconds when the game started; it survives server crashes
// because it travels with the lobby state.
// ElapsedMs is the final game duration in milliseconds, set only when the game ends.
type GameInfo struct {
	Snippet      string
	CorrectChars int64
	WrongChars   int64
	Revision     int64
	StartTimeMs  int64
	ElapsedMs    int64
}

// NewGameInfo creates a new GameInfo for a game starting with the given snippet.
// CorrectChars, WrongChars and Revision are initialised to zero.
func NewGameInfo(snippet string) *GameInfo {
	return &GameInfo{Snippet: snippet}
}
