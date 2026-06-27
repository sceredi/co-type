// Package game_messages defines the messages for the game page.
package game_messages

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/common/domain"
)

// GameEndMsg is a message that indicates that a game has ended.
type GameEndMsg struct {
	Stats domain.GameStats
}

// NewGameEndCmd returns a command that creates a GameEndMsg.
func NewGameEndCmd(stats domain.GameStats) tea.Cmd {
	return func() tea.Msg {
		return GameEndMsg{Stats: stats}
	}
}

// KeyPressResultMsg carries the lobby state returned after a key press, or an error.
type KeyPressResultMsg struct {
	Lobby *domain.Lobby
	Err   error
}

// NewSendKeyPressCmd returns a Cmd that sends a key press to the server and emits a KeyPressResultMsg.
func NewSendKeyPressCmd(lobbySvc service.LobbyService, lobbyID string, player *domain.Player, key string, isBackspace bool) tea.Cmd {
	return func() tea.Msg {
		l, err := lobbySvc.SendKeyPress(lobbyID, player, key, isBackspace)
		return KeyPressResultMsg{Lobby: l, Err: err}
	}
}
