package header

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
)

var (
	headerCellStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.NormalBorder())
	firstCellStyle = headerCellStyle.BorderLeft(true)
)

type Model struct {
	gameName    string
	inviteCode  string
	playerCount int
	maxPlayers  int
	width       int
}

func New(gameName, inviteCode string, maxPlayers int) Model {
	return Model{
		gameName:   gameName,
		inviteCode: inviteCode,
		maxPlayers: maxPlayers,
		width:      80,
	}
}

func (m Model) SetWidth(w int) Model {
	m.width = w
	return m
}

func (m Model) SetPlayerCount(n int) Model {
	m.playerCount = n
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	cells := make([]string, 4)
	cells[0] = firstCellStyle.Render(styles.NewLabelBold("co-type"))
	cells[1] = headerCellStyle.Render(m.gameName)
	cells[2] = headerCellStyle.Render(fmt.Sprintf("[ %s ] %s", styles.NewLabelBold("i"), m.inviteCode))
	cells[3] = headerCellStyle.Render(fmt.Sprintf("%d/%d", m.playerCount, m.maxPlayers))
	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}
