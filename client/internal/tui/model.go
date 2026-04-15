package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby"
)

type page int

const (
	lobbyPage page = iota
	gamePage
)

type Model struct {
	width       int
	height      int
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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

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
	if m.width > 0 && m.height > 0 {
		v = lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			v,
		)
	}
	view := tea.NewView(v)
	view.AltScreen = true
	return view
}
