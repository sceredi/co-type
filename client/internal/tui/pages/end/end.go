// Package end defines the end page of the TUI.
package end

import (
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

// Model is the model for the end page.
type Model struct {
	stats domain.GameStats
}

// New creates a new and page model.
func New(stats domain.GameStats) Model {
	return Model{stats: stats}
}

// Init initializes the end page model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the end page model based on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	log.Printf("Received %v", msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the end page.
func (m Model) View() string {
	header := "Game statistics:"
	d := m.stats.TotalTime
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	duration := fmt.Sprintf("Duration -> %02d:%02d (mm:ss)", minutes, seconds)
	cps := fmt.Sprintf("Characters per second -> %.02f", float64(len(m.stats.Lobby.Snippet))/m.stats.TotalTime.Seconds())
	v := lipgloss.JoinVertical(lipgloss.Left, duration, cps)
	v = lipgloss.JoinVertical(lipgloss.Center, header, v, "\nUse \"ctrl+c\" to close co-type")
	return styles.NewContainer(v)
}
