package settings

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

type focusSlot int

const (
	focusAllowed focusSlot = iota
	focusBlocked
	focusBackspace
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

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var c tea.Cursor
	if !m.allowedList.VirtualCursor() {
		c = *m.allowedList.Cursor()
		c.Y += lipgloss.Height(allowedDescription)
	}
	v := lipgloss.JoinVertical(
		lipgloss.Center,
		allowedDescription,
		m.allowedList.View(),
		blockedDescription,
		m.blockedList.View(),
	)
	return styles.NewContainer(v)
}
