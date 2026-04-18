package settings

import (
	"image/color"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

type ChangeFocusMsg struct {
	IsFocussed bool
}

type focusSlot int

const (
	focusAllowed focusSlot = iota
	focusBlocked
)

const (
	allowedDescription = "Allowed characters regex"
	blockedDescription = "Blocked characters regex"
)

type Model struct {
	isFocussed  bool
	focus       focusSlot
	allowedList textinput.Model
	blockedList textinput.Model
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

func New(allowed string, blocked string) Model {
	allowedList := newTextinput("eg. /[a-z]/", allowed)
	blockedList := newTextinput("eg. /[0-9]/", blocked)
	return Model{
		isFocussed:  false,
		focus:       focusAllowed,
		allowedList: allowedList,
		blockedList: blockedList,
	}
}

func (m *Model) unfocus() {
	m.allowedList.Blur()
	m.blockedList.Blur()
}
func (m *Model) updateFocus() {
	m.unfocus()
	if m.isFocussed {
		switch m.focus {
		case focusAllowed:
			m.allowedList.Focus()
		case focusBlocked:
			m.blockedList.Focus()
		}
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
	case ChangeFocusMsg:
		m.isFocussed = msg.IsFocussed
		m.updateFocus()
		m.allowedList.Focus()
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
			m.focus = m.focus % (focusBlocked + 1)
			m.updateFocus()
		}
	}
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var color color.Color
	if m.isFocussed {
		color = styles.Blue
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
	}
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		allowedDescription,
		m.allowedList.View(),
		blockedDescription,
		m.blockedList.View(),
	)
	return styles.NewContainer(v, color)
}
