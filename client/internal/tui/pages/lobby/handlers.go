package lobby

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/settings"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
	"github.com/sceredi/co-type/common/domain"
)

func (m *Model) handleReadyCmd() tea.Cmd {
	m.player.IsReady = !m.player.IsReady
	shouldStart := true
	for _, p := range m.lobby.Players {
		if !p.IsReady {
			shouldStart = false
			break
		}
	}
	if shouldStart {
		// TODO: hardcoded value
		m.lobby.Snippet = strings.Trim(`
m.player.IsReady = !m.player.IsReady
shouldStart := true
for _, p := range m.lobby.Players {
	if !p.IsReady {
		shouldStart = false
		break
	}
}`, "\r\n\t")
		game := domain.Game{
			Lobby: *m.lobby,
			State: domain.GameState{
				CorrectChars: 15,
				LastChar:     17,
				Status:       domain.Running,
			},
		}
		return lobby_messages.NewStartGameCmd(game)
	}
	return nil
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
