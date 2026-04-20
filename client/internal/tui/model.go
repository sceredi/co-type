// Package tui defines the TUI model, which is the main component of the TUI.
// It manages the state of the TUI and the movement between pages.
package tui

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby"
	"github.com/sceredi/co-type/common/domain"
)

type page int

const (
	lobbyPage page = iota
	gamePage
)

// Model is the main model of the TUI.
type Model struct {
	width       int
	height      int
	currentPage page
	lobby       lobby.Model
}

// New creates a new TUI model.
func New() Model {
	// TODO: remove hardcoded values
	p := domain.NewPlayer("sceredi")
	p.AllowedCharacters = "/[a-z]/"
	l := domain.NewLobby("lobby1", &p)
	p2 := domain.NewPlayer("player2")
	p3 := domain.NewPlayer("player3")
	p4 := domain.NewPlayer("player4")
	l.AddPlayers(&p2, &p3, &p4)
	return Model{
		lobby: lobby.New(
			&p,
			&l,
		),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	cmd := m.lobby.Init()
	return cmd
}

// Update updates the model based on the given message and returns the updated model and any commands to execute.
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

// View returns the view of the model.
func (m Model) View() tea.View {
	var v string
	switch m.currentPage {
	case lobbyPage:
		v = m.lobby.View()
	default:
		// TODO: not yet implemented
		log.Fatal("not yet implemented")
	}
	if m.width < pages.Width || m.height < pages.Height {
		v = "Your terminal is too small to display the game. \nPlease increase the size of your terminal or \ndecrease the font size [ctrl + minus]."
	}
	v = lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		v,
	)
	view := tea.NewView(v)
	view.AltScreen = true
	return view
}
