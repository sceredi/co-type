// Package tui defines the TUI model, which is the main component of the TUI.
// It manages the state of the TUI and the movement between pages.
package tui

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	"github.com/sceredi/co-type/client/internal/tui/pages/end"
	"github.com/sceredi/co-type/client/internal/tui/pages/game"
	game_messages "github.com/sceredi/co-type/client/internal/tui/pages/game/messages"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
	"github.com/sceredi/co-type/client/internal/tui/pages/welcome"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
)

type page int

const (
	welcomePage page = iota
	lobbyPage
	gamePage
	endPage
)

// Model is the main model of the TUI.
type Model struct {
	width  int
	height int
	page   page

	welcome welcome.Model
	lobby   lobby.Model
	game    game.Model
	end     end.Model
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
	}

	switch m.page {
	case welcomePage:
		switch msg := msg.(type) {
		case welcome_messages.JoinLobbyMsg:
			m.lobby = lobby.New(msg.Lobby.Host, &msg.Lobby)
			m.page = lobbyPage
		}
		m.welcome, cmd = m.welcome.Update(msg)
	case lobbyPage:
		switch msg := msg.(type) {
		case lobby_messages.LeaveMsg:
			m.page = welcomePage
		case lobby_messages.StartGameMsg:
			m.game = game.New(msg.Game)
			m.page = gamePage
		}
		m.lobby, cmd = m.lobby.Update(msg)
	case gamePage:
		switch msg := msg.(type) {
		case game_messages.GameEndMsg:
			m.end = end.New(msg.Stats)
			m.page = endPage
		}
		m.game, cmd = m.game.Update(msg)
	case endPage:
		m.end, cmd = m.end.Update(msg)

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
	case gamePage:
		v = m.game.View()
	case endPage:
		v = m.end.View()
	default:
		log.Fatal("model.go -> Unexpected page value in View()")
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
