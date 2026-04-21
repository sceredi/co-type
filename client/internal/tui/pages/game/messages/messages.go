// Package game_messages defines the messages for the game page.
package game_messages

import (
	tea "charm.land/bubbletea/v2"
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
