package lobby

import "github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/settings"

func (m *Model) handleReadyCmd() {
	m.player.IsReady = !m.player.IsReady
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
