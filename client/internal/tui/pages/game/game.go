// Package game defines the game page of the TUI.
package game

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

// Model is the model for the game page.
type Model struct {
	width  int
	height int

	game domain.Game
}

// New creates a new game model.
func New(game domain.Game) Model {
	return Model{game: game}
}

// Init initializes the game model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the game model base on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		log.Printf("game: pressed %s", msg.String())
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the game page.
func (m Model) View() string {
	style := lipgloss.NewStyle()
	ce := m.game.State.CorrectChars
	end := m.game.State.LastChar
	correct := style.Background(styles.DarkGreen).Render(m.game.Lobby.Snippet[:ce])
	wrong := style.Background(styles.DarkRed).Render(m.game.Lobby.Snippet[ce:end])
	rest := m.game.Lobby.Snippet[end:]
	v := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Left).
		Render(correct + wrong + rest)
	return styles.NewContainer(v)
}
