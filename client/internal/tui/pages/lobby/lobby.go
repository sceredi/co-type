package lobby

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/header"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/playerlist"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/settings"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

type focusSlot int

const (
	focusPlayersList focusSlot = iota
	focusSettings
)

// The Model struct represents the state of the lobby page.
type Model struct {
	focus          focusSlot
	selectedPlayer int
	player         *domain.Player
	lobby          *domain.Lobby

	settings settings.Model
}

func New(player *domain.Player, lobby *domain.Lobby) Model {
	m := Model{
		focus:          focusPlayersList,
		selectedPlayer: 0,
		player:         player,
		lobby:          lobby,
	}
	p := m.lobby.Players[m.selectedPlayer]
	allowed := p.AllowedCharacters
	blocked := p.BlockedCharacters
	m.settings = settings.New(allowed, blocked)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "right":
			m.focus = (m.focus + 1) % (focusSettings + 1)
			m.settings, _ = m.settings.Update(settings.ChangeFocusMsg{IsFocussed: m.focus == focusSettings})
		}
	}

	if m.focus == focusPlayersList {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up":
				handlePlayerSelect(m.selectedPlayer-1, &m)
			case "down":
				handlePlayerSelect(m.selectedPlayer+1, &m)
			case "r":
				handleReadyCmd(&m)
			}
		}
	} else {
		m.settings, cmd = m.settings.Update(msg)
	}

	cmds = append(cmds, cmd)
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

	body := lipgloss.JoinHorizontal(lipgloss.Center,
		playerlist.Render(m.selectedPlayer, m.lobby, m.focus == focusPlayersList),
		m.settings.View(),
	)

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, readyBtn, leaveBtn)

	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header.Render(m.player, m.lobby),
		body,
		buttons,
	)
	return v
}
