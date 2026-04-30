// Package welcome defines the welcome page of the TUI.
package welcome

import (
	"log"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/client/internal/tui/pages/common"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

type focusSlot int

const (
	focusName focusSlot = iota
	focusCode
	focusCreate
	focusJoin
)

// Model is the model for the welcome page.
type Model struct {
	focus focusSlot

	nameTi textinput.Model
	codeTi textinput.Model

	error string

	ds service.DiscoveryService
}

// New creates a new welcome model.
func New(ds service.DiscoveryService) Model {
	tn := common.NewTextinput("Choose your username", "")
	tn.Focus()
	tc := common.NewTextinput("CT2026", "")
	tc.CharLimit = 10
	return Model{
		focus: focusName,

		nameTi: tn,
		codeTi: tc,

		ds: ds,
	}
}

// Init initializes the welcome model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) updateFocus() {
	m.nameTi.Blur()
	m.codeTi.Blur()
	switch m.focus {
	case focusName:
		m.nameTi.Focus()
	case focusCode:
		m.codeTi.Focus()
	default:
	}
}

// Update updates the welcome model based on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			v := msg.String()
			if v == "tab" || v == "down" {
				m.focus++
			} else {
				m.focus--
			}
			m.focus = (m.focus + focusJoin + 1) % (focusJoin + 1)
			m.updateFocus()
		case "enter":
			switch m.focus {
			case focusName, focusCode:
				m.focus++
				m.updateFocus()
			case focusCreate:
				// TODO: create new lobby
				cmds = append(cmds, welcome_messages.NewCreateLobbyCmd(m.ds, m.codeTi.Value(), m.nameTi.Value()))
			case focusJoin:
				// TODO: join the lobby
				cmds = append(cmds, welcome_messages.NewJoinLobbyCmd(m.codeTi.Value(), m.nameTi.Value()))
			default:
				log.Fatalf("welcome: unexpected focus value %d", m.focus)
			}
		}
	case welcome_messages.JoinLobbyErrorMsg:
		m.error = msg.Error
	}
	switch m.focus {
	case focusName:
		m.nameTi, cmd = m.nameTi.Update(msg)
	case focusCode:
		m.codeTi, cmd = m.codeTi.Update(msg)
	default:
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

const keybinds = "↑/↓/tab/s+tab: select field\nenter: confirm • ctrl+c: quit"

// View renders the welcome page.
func (m Model) View() string {
	header := styles.LabelBold.PaddingBottom(1).Render("CO-TYPE")
	var content string
	createButtonStyle := styles.ButtonDefault
	joinButtonStyle := styles.ButtonDefault
	switch m.focus {
	case focusCreate:
		createButtonStyle = styles.ButtonBlue
	case focusJoin:
		joinButtonStyle = styles.ButtonBlue
	default:
	}
	createBtn := createButtonStyle.Render("Create lobby")
	joinBtn := joinButtonStyle.Render("Join lobby")
	content = lipgloss.JoinVertical(
		lipgloss.Center,
		styles.PaddingVertical.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Username:",
				m.nameTi.View(),
			),
		),
		styles.PaddingVertical.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Lobby code:",
				m.codeTi.View(),
			),
		),
		createBtn,
		joinBtn,
	)
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		content,
		keybinds,
	)
	if m.error != "" {
		v = lipgloss.JoinVertical(
			lipgloss.Center,
			v,
			lipgloss.NewStyle().Foreground(styles.Red).Render(m.error),
		)
	}
	return styles.NewContainer(v)
}
