package lobby

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	header "github.com/sceredi/co-type/client/internal/tui/pages/lobby/components"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

// --- Messages ---

type countdownTickMsg struct{}

func countdownTick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return countdownTickMsg{}
	})
}

func manageCountDownTickMsg(m *Model, cmds []tea.Cmd) []tea.Cmd {
	if m.countdown > 0 {
		m.countdown--
		m.statusMsg = fmt.Sprintf("Game starting in %d...", m.countdown)
		cmds = append(cmds, countdownTick())
	} else if m.countdown == 0 {
		m.statusMsg = "Launching game!"
		m.countdown = -1
	}
	return cmds
}

// --- Focus ---

// The focusSlot type represents the different focusable areas in the lobby page.
type focusSlot int

const (
	focusSettings focusSlot = iota
	focusReady
	focusLeave
)

// The Model struct represents the state of the lobby page.
type Model struct {
	header header.Model // TODO: add remaining components

	isReady   bool
	focus     focusSlot
	countdown int // -1 if not counting
	statusMsg string
	width     int
	height    int
}

func New() Model {
	m := Model{
		header:    header.New("CO-TYPE", "CT-2DMK", 8),
		countdown: -1,
		statusMsg: "Use tab to switch panels - [r] ready up - [q] quit",
		focus:     focusReady,
	}
	m.header = m.header.SetPlayerCount(1)
	return m
}

func (m Model) SetSize(w, h int) Model {
	m.width = w
	m.height = h
	m.header = m.header.SetWidth(w)
	return m
}

func (m Model) Init() tea.Cmd {
	cmd := m.header.Init()
	return cmd
}

func (m Model) setFocus(slot focusSlot) Model {
	m.focus = slot
	// TODO: after adding components blur everything and unblur the selected one
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case countdownTickMsg:
		cmds = manageCountDownTickMsg(&m, cmds)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			next := (int(m.focus) + 1) % int(focusReady)
			m = m.setFocus(focusSlot(next))

		case "r":
			if m.countdown == -1 {
				m.isReady = !m.isReady
				if m.isReady {
					m.statusMsg = "You are ready!"
					cmds = append(cmds, countdownTick()) // TODO: Remove
					// TODO: Check if all ready then start countdown
				} else {
					m.statusMsg = "Waiting... press [r] to ready up"
				}
			}
		}
	}

	// --- Delegate to children ---
	m.header, cmd = m.header.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var readyLabel string
	if m.isReady {
		readyLabel = "[ r ] Ready!"
	} else {
		readyLabel = "[ r ] Ready Up"
	}

	readyBtn := styles.ButtonPrimary.Render(readyLabel)
	leaveBtn := styles.ButtonDanger.Render("[ esc ] Leave")

	buttons := lipgloss.JoinHorizontal(lipgloss.Right, readyBtn, leaveBtn)

	v := lipgloss.JoinVertical(
		lipgloss.Center,
		m.header.View(),
		buttons,
	)
	return v
}
