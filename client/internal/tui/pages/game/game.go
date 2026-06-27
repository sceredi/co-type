// Package game defines the game page of the TUI.
package game

import (
	"context"
	"log/slog"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	game_messages "github.com/sceredi/co-type/client/internal/tui/pages/game/messages"
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

// Update updates the game model based on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			// TODO: this is hardcoded
			tt := time.Now().Round(0).Add(-(120) * time.Second)
			td := time.Since(tt)
			cmd = game_messages.NewGameEndCmd(domain.GameStats{
				TotalTime: td,
				Lobby:     m.game.Lobby,
			})
			cmds = append(cmds, cmd)
		}

		slog.DebugContext(context.Background(), "key pressed",
			slog.String("key", msg.String()),
		)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) renderPlayers() string {
	rows := make([]string, len(m.game.Lobby.Players)+1)
	rows[0] = "Players:"
	for i, p := range m.game.Lobby.Players {
		rows[i+1] = p.Name
	}
	v := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return v
}

// View renders the game page.
func (m Model) View() string {
	style := lipgloss.NewStyle()
	ce := m.game.State.CorrectChars
	end := m.game.State.LastChar
	correct := style.Background(styles.DarkGreen).Render(m.game.Lobby.Snippet[:ce])
	wrong := style.Background(styles.DarkRed).Render(m.game.Lobby.Snippet[ce:end])
	rest := m.game.Lobby.Snippet[end:]
	board := lipgloss.NewStyle().
		Width(max(m.width, pages.Width) * 4 / 5).
		Height(max(m.height, pages.Height) - 4).
		Align(lipgloss.Left).
		PaddingRight(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		Render(correct + wrong + rest)
	players := lipgloss.NewStyle().PaddingLeft(2).Render(m.renderPlayers())

	v := lipgloss.JoinHorizontal(lipgloss.Top, board, players)

	v = styles.NewContainer(v)

	if m.game.State.Status == domain.Paused {
		base := lipgloss.NewLayer(v)
		topContent := styles.NewContainer(
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Game paused",
				"Press [ r ] to resume",
			),
		)
		x := (lipgloss.Width(v) / 2) - (lipgloss.Width(topContent) / 2)
		y := (lipgloss.Height(v) / 2) - (lipgloss.Height(topContent) / 2)
		top := lipgloss.NewLayer(
			topContent,
		).X(x).Y(y).Z(1)
		comp := lipgloss.NewCompositor(base, top)
		return comp.Render()
	}
	return v
}
