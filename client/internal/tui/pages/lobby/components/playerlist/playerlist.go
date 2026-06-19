// Package playerlist defines the player list component of the lobby page.
package playerlist

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/sceredi/co-type/common/domain"
)

// Render renders the player list component for the given selected player and lobby.
func Render(selectedPlayer int, lobby *domain.Lobby) string {
	rows := make([][]string, len(lobby.Players))
	for i, player := range lobby.Players {
		pointer := " "
		if i == selectedPlayer {
			pointer = "→"
		}
		rows[i] = append(rows[i], pointer)
		rows[i] = append(rows[i], player.Name)
		ready := " "
		if player.IsReady {
			ready = "r"
		}
		rows[i] = append(rows[i], ready)
	}
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		StyleFunc(func(_, _ int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(1)
		}).
		Rows(rows...)
	return t.Render()
}
