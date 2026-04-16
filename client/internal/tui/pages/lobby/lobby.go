package lobby

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/header"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/playerlist"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

// The Model struct represents the state of the lobby page.
type Model struct {
	selectedPlayer int
	player         *domain.Player
	lobby          *domain.Lobby
}

func New(player *domain.Player, lobby *domain.Lobby) Model {
	m := Model{
		selectedPlayer: 0,
		player:         player,
		lobby:          lobby,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			handlePlayerSelect(m.selectedPlayer-1, &m)
		case "down", "j":
			handlePlayerSelect(m.selectedPlayer+1, &m)
		case "r":
			handleReadyCmd(&m)
		}

	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var readyLabel string
	if m.player.IsReady {
		readyLabel = "[ r ] Ready!"
	} else {
		readyLabel = "[ r ] Ready Up"
	}

	readyBtn := styles.ButtonPrimary.Render(readyLabel)
	leaveBtn := styles.ButtonDanger.Render("[ esc ] Leave")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, readyBtn, leaveBtn)

	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header.Render(m.player, m.lobby),
		playerlist.Render(m.selectedPlayer, m.lobby),
		buttons,
	)
	return v
}
