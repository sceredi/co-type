// Package joinlobby provides a bubbletea component to join a lobby by entering a lobby code.
package joinlobby

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/common"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

// Model is the bubbletea model for the join lobby component.
type Model struct {
	lobbyCode textinput.Model
	error     string
}

// New creates a new join lobby model.
func New() Model {
	ti := common.NewTextinput("CT20KM", "")
	ti.Focus()
	return Model{lobbyCode: ti}
}

// Init initializes the join lobby model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the join lobby model based on the given message
// and returns the updated model and any commands to execute.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			// TODO: join the given lobby
			m.error = ""
			return m, welcome_messages.NewJoinLobbyCmd(m.lobbyCode.Value())
		case "esc":
			return m, welcome_messages.NewLeaveJoinLobbyCmd()
		}
	case welcome_messages.JoinLobbyErrorMsg:
		m.error = msg.Error
	}
	m.lobbyCode, cmd = m.lobbyCode.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the join lobby component.
func (m Model) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"Lobby code:",
		m.lobbyCode.View(),
	)
	btnJoin := styles.ButtonPrimary.Render("[ j ] Join")
	btnLeave := styles.ButtonDanger.Render("[ esc ] Leave")
	buttons := lipgloss.JoinHorizontal(lipgloss.Center, btnJoin, btnLeave)
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		content,
		buttons,
	)
	if m.error != "" {
		v = lipgloss.JoinVertical(
			lipgloss.Center,
			v,
			lipgloss.NewStyle().Foreground(styles.Red).Render(m.error),
		)
	}
	return v
}
