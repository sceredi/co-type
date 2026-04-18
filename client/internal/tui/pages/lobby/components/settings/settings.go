package settings

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

type focusSlot int

const (
	focusAllowed focusSlot = iota
	focusBlocked
	focusConfim
	focusCancel
)

const (
	allowedDescription = "Allowed characters regex"
	blockedDescription = "Blocked characters regex"
)

type Model struct {
	focus            focusSlot
	player           *domain.Player
	allowedList      textinput.Model
	blockedList      textinput.Model
	backspaceAllowed bool
}

func newTextinput(placeholder string, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.SetWidth(20)
	suggestions := []string{
		"/[a-z]/",
		"/[0-9]/",
		"/[a-zA-Z0-9]/",
	}
	ti.SetSuggestions(suggestions)
	ti.SetValue(value)
	return ti
}

func New(player *domain.Player) Model {
	allowedList := newTextinput("eg. /[a-z]/", player.AllowedCharacters)
	blockedList := newTextinput("eg. /[0-9]/", player.BlockedCharacters)
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
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.focus {
	case focusAllowed:
		m.allowedList, cmd = m.allowedList.Update(msg)
	case focusBlocked:
		m.blockedList, cmd = m.blockedList.Update(msg)
	}
	return cmd
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			v := msg.String()
			if v == "tab" || v == "enter" || v == "down" {
				m.focus++
			} else {
				m.focus--
			}
			m.focus = m.focus % (focusCancel + 1)
			m.updateFocus()
		}
	}
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := styles.PaddingVertical.Render(fmt.Sprintf("Editing settings for %s", styles.NewLabelBold(m.player.Name)))
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
	confirmButtonStyle := styles.ButtonDefault
	cancelButtonStyle := styles.ButtonDefault
	switch m.focus {
	case focusConfim:
		confirmButtonStyle = styles.ButtonPrimary
	case focusCancel:
		cancelButtonStyle = styles.ButtonDanger
	}
	confirmBtn := confirmButtonStyle.Render("Confirm")
	cancelBtn := cancelButtonStyle.Render("Cancel")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, confirmBtn, cancelBtn)
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		allowedDescription,
		m.allowedList.View(),
		blockedDescription,
		m.blockedList.View(),
		buttons,
	)
	return styles.NewContainer(v)
}
