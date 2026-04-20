// Package tui defines the TUI model, which is the main component of the TUI.
// It manages the state of the TUI and the movement between pages.
package tui

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby"
	"github.com/sceredi/co-type/client/internal/tui/pages/welcome"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
)

type page int

const (
	welcomePage page = iota
	lobbyPage
	gamePage
)

// Model is the main model of the TUI.
type Model struct {
	width  int
	height int
	page   page

	welcome welcome.Model
	lobby   lobby.Model
}

// New creates a new TUI model.
func New() Model {
	return Model{
		page: welcomePage,

		welcome: welcome.New(),
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the model based on the given message and returns the updated model and any commands to execute.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case welcome_messages.JoinLobbyMsg:
		m.lobby = lobby.New(msg.Lobby.Host, &msg.Lobby)
		m.page = lobbyPage
	}

	switch m.page {
	case welcomePage:
		m.welcome, cmd = m.welcome.Update(msg)
	case lobbyPage:
		m.lobby, cmd = m.lobby.Update(msg)
	default:
		log.Fatalf("model: ERROR Unexpected page %d", m.page)
	}

	return m, cmd
}

// View returns the view of the model.
func (m Model) View() tea.View {
	var v string
	switch m.page {
	case welcomePage:
		v = m.welcome.View()
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
