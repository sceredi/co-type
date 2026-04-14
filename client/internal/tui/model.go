package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby"
)

type page int

const (
	lobbyPage page = iota
	gamePage
)

type Model struct {
	currentPage page
	lobby       lobby.Model
}

func New() Model {
	return Model{
		lobby: lobby.New(),
	}
}

func (m Model) Init() tea.Cmd {
	cmd := m.lobby.Init()
	return cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.lobby, cmd = m.lobby.Update(msg)

	// TODO: Manage movement between pages

	// TODO: Update only the current page
	return m, cmd
}

func (m Model) View() tea.View {
	var v string
	switch m.currentPage {
	case lobbyPage:
		v = m.lobby.View()
	}
	view := tea.NewView(v)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}
