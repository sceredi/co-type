package lobby

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/header"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

// --- Messages ---
type PlayerReadyMsg struct{}

func newPlayerReadyMsg() tea.Cmd {
	return func() tea.Msg {
		return PlayerReadyMsg{}
	}
}

// The Model struct represents the state of the lobby page.
type Model struct {
	player *domain.Player
	lobby  *domain.Lobby
}

func New(player *domain.Player, lobby *domain.Lobby) Model {
	m := Model{
		player: player,
		lobby:  lobby,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case PlayerReadyMsg:
		log.Printf("Updating ready from %t\n", m.player.IsReady)
		m.player.IsReady = !m.player.IsReady

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			log.Println("Pressed ready")
			cmds = append(cmds, newPlayerReadyMsg())
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

	buttons := lipgloss.JoinHorizontal(lipgloss.Right, readyBtn, leaveBtn)

	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header.Render(m.player, m.lobby),
		buttons,
	)
	return v
}
