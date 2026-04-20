// Package welcome defines the welcome page of the TUI.
package welcome

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/welcome/components/joinlobby"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

type focusSlot int

const (
	focusWelcome focusSlot = iota
	focusJoin
)

// Model is the model for the welcome page.
type Model struct {
	focus     focusSlot
	joinlobby joinlobby.Model
}

// New creates a new welcome model.
func New() Model {
	return Model{
		focus:     focusWelcome,
		joinlobby: joinlobby.New(),
	}
}

// Init initializes the welcome model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the welcome model based on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	if m.focus == focusWelcome {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "c":
				// TODO: create new lobby
			case "j":
				// TODO: join the lobby
				m.focus = focusJoin
			}
		}
	} else {
		switch msg.(type) {
		case welcome_messages.LeaveJoinLobbyMsg:
			m.focus = focusWelcome
		}
		m.joinlobby, cmd = m.joinlobby.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the welcome page.
func (m Model) View() string {
	header := styles.LabelBold.PaddingBottom(1).Render("CO-TYPE")
	var content string
	if m.focus == focusWelcome {
		createBtn := styles.ButtonDefault.Render(fmt.Sprintf("[ %s ] Create lobby", styles.NewLabelBold("c")))
		joinBtn := styles.ButtonDefault.Render(fmt.Sprintf("[ %s ] Join lobby", styles.NewLabelBold("j")))
		content = lipgloss.JoinVertical(lipgloss.Center, createBtn, joinBtn)
	} else {
		content = m.joinlobby.View()
	}
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		content,
	)
	return styles.NewContainer(v)
}
