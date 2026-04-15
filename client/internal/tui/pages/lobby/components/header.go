package header

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
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
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderHeader(false).
		BorderRow(false).
		Headers(
			styles.NewLabelBold("co-type"),
			m.gameName,
			fmt.Sprintf("[ %s ] %s", styles.NewLabelBold("i"), m.inviteCode),
			fmt.Sprintf("%d/%d", m.playerCount, m.maxPlayers),
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Width(pages.Width/4 - 1).Align(lipgloss.Center)
		})

	return t.Render()
}
