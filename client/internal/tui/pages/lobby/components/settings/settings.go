// Package settings provides a bubbletea component to edit player settings
package settings

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/common"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

type focusSlot int

const (
	focusAllowed focusSlot = iota
	focusBlocked
	focusBackspace
	focusConfim
	focusCancel
)

const (
	allowedDescription = "Allowed characters regex"
	blockedDescription = "Blocked characters regex"
)

// Model is the bubbletea model for the settings component.
type Model struct {
	focus            focusSlot
	player           *domain.Player
	allowedList      textinput.Model
	blockedList      textinput.Model
	backspaceAllowed bool
}

// New creates a new settings model for the given player.
func New(player *domain.Player) Model {
	allowedList := common.NewTextinput("eg. /[a-z]/", player.AllowedCharacters)
	blockedList := common.NewTextinput("eg. /[0-9]/", player.BlockedCharacters)
	allowedList.Focus()
	return Model{
		focus:            focusAllowed,
		player:           player,
		allowedList:      allowedList,
		blockedList:      blockedList,
		backspaceAllowed: player.CanDelete,
	}
}

func (m *Model) updateFocus() {
	m.allowedList.Blur()
	m.blockedList.Blur()
	switch m.focus {
	case focusAllowed:
		m.allowedList.Focus()
	case focusBlocked:
		m.blockedList.Focus()
	default:
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.focus {
	case focusAllowed:
		m.allowedList, cmd = m.allowedList.Update(msg)
	case focusBlocked:
		m.blockedList, cmd = m.blockedList.Update(msg)
	default:
	}
	return cmd
}

// Init initializes the model. It starts the blinking cursor for the text inputs.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) confirmSaveCmd() tea.Cmd {
	return lobby_messages.NewUpdatePlayerCmd(
		m.player,
		m.allowedList.Value(),
		m.blockedList.Value(),
		m.backspaceAllowed,
	)
}

func (m *Model) cancelSaveCmd() tea.Cmd {
	return lobby_messages.NewUpdatePlayerCmd(
		m.player,
		m.player.AllowedCharacters,
		m.player.BlockedCharacters,
		m.player.CanDelete,
	)
}

const (
	spacek = "space"
	tabk   = "tab"
	stabk  = "shift+tab"
	enterk = "enter"
	upk    = "up"
	downk  = "down"
)

// Update updates the model based on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case spacek:
			if m.focus == focusBackspace {
				m.backspaceAllowed = !m.backspaceAllowed
			}
		case tabk, stabk, enterk, upk, downk:
			if msg.String() == enterk && m.focus == focusConfim {
				cmd := m.confirmSaveCmd()
				return m, tea.Batch(cmd, lobby_messages.NewCloseSettingsCmd())
			}
			if msg.String() == enterk && m.focus == focusCancel {
				cmd := m.cancelSaveCmd()
				return m, tea.Batch(cmd, lobby_messages.NewCloseSettingsCmd())
			}
			if msg.String() == enterk && m.focus == focusBackspace {
				m.backspaceAllowed = !m.backspaceAllowed
			}
			v := msg.String()
			if v == tabk || v == enterk || v == downk {
				m.focus++
			} else {
				m.focus--
			}
			m.focus = (m.focus + focusCancel + 1) % (focusCancel + 1)
			m.updateFocus()
		}
	}
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

const keybinds = "↑/↓/tab/s+tab: select field\nspace: toggle • enter: confirm\nctrl+c: quit"

// View renders the settings editor.
func (m Model) View() string {
	header := styles.PaddingVertical.Render(
		fmt.Sprintf("Editing settings for %s", styles.NewLabelBold(m.player.Name)),
	)
	var c tea.Cursor
	var focussedList textinput.Model
	if m.allowedList.Focused() {
		focussedList = m.allowedList
	} else {
		focussedList = m.blockedList
	}
	if !focussedList.VirtualCursor() {
		c = *m.allowedList.Cursor()
		c.Y += lipgloss.Height(allowedDescription)
	}

	allowedView := styles.PaddingVertical.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			allowedDescription,
			m.allowedList.View(),
		),
	)

	blockedView := styles.PaddingVertical.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			blockedDescription,
			m.blockedList.View(),
		),
	)
	backspaceStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false)
	backspaceText := "enabled"
	backspaceTextStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if !m.backspaceAllowed {
		backspaceText = "disabled"
		backspaceTextStyle = backspaceTextStyle.Foreground(styles.Red)
	}
	if m.focus == focusBackspace {
		backspaceStyle = backspaceStyle.BorderForeground(styles.Blue)
	}
	backspace := fmt.Sprintf("Backspace: %s", backspaceTextStyle.Render(backspaceText))
	backspace = backspaceStyle.Render(backspace)

	editorView := lipgloss.JoinVertical(
		lipgloss.Left,
		allowedView,
		blockedView,
		backspace,
	)

	confirmButtonStyle := styles.ButtonDefault
	cancelButtonStyle := styles.ButtonDefault
	switch m.focus {
	case focusConfim:
		confirmButtonStyle = styles.ButtonPrimary
	case focusCancel:
		cancelButtonStyle = styles.ButtonDanger
	default:
	}
	confirmBtn := confirmButtonStyle.Render("Confirm")
	cancelBtn := cancelButtonStyle.Render("Cancel")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, confirmBtn, cancelBtn)
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		editorView,
		buttons,
		keybinds,
	)
	return styles.NewContainer(v)
}
