package playerlist

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

func Render(selectedPlayer int, lobby *domain.Lobby) string {
	rows := make([][]string, len(lobby.Players))
	for i, player := range lobby.Players {
		pointer := " "
		if i == selectedPlayer {
			pointer = "→"
		}
		rows[i] = append(rows[i], pointer)
		name := player.Name
		if player.Name == lobby.Host.Name {
			name = styles.NewLabelBold(player.Name)
		}
		rows[i] = append(rows[i], name)
		ready := " "
		if player.IsReady {
			ready = "r"
		}
		rows[i] = append(rows[i], ready)
	}
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderColumn(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(1)
		}).
		Rows(rows...)
	return t.Render()
}
