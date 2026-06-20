package lobby

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/settings"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
)

func (m *Model) handleReadyCmd() tea.Cmd {
	return lobby_messages.NewReadyCmd(m.ls, m.lobby.ID, m.player.Name)
}

func (m *Model) handlePlayerSelect(selectedPlayer int) {
	if selectedPlayer < 0 {
		m.selectedPlayer = len(m.lobby.Players) - 1
	} else {
		m.selectedPlayer = selectedPlayer % len(m.lobby.Players)
	}
}

func (m *Model) handleOpenSettingsModal() {
	m.focus = focusSettings
	m.settings = settings.New(m.lobby.Players[m.selectedPlayer])
}
