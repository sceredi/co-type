package header

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

func Render(player *domain.Player, lobby *domain.Lobby) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(true).
		BorderHeader(false).
		BorderRow(false).
		Headers(
			styles.NewLabelBold("co-type"),
			fmt.Sprintf("[ %s ] %s", styles.NewLabelBold("i"), lobby.ID),
			styles.NewLabelBold(player.Name),
			fmt.Sprintf("%d/%d", len(lobby.Players), lobby.MaxPlayers),
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Width(pages.Width / 4).Align(lipgloss.Center)
		})

	return t.Render()
}
