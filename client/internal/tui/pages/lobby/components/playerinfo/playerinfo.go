// Package playerinfo defines the component that renders the player's information in the lobby page.
package playerinfo

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

const (
	allowedDesc   = "Allowed characters"
	blockedDesc   = "Blocked characters"
	backspaceDesc = "Backspace:"
)

// Render renders the player information component for the given player.
func Render(player domain.Player) string {
	renderer := lipgloss.NewStyle().Padding(1)
	var allowed, blocked string
	if player.AllowedCharacters != "" {
		allowed = fmt.Sprintf("%s: \"%s\"", allowedDesc, player.AllowedCharacters)
	} else {
		allowed = fmt.Sprintf("%s not yet set", allowedDesc)
	}
	if player.BlockedCharacters != "" {
		blocked = fmt.Sprintf("%s: \"%s\"", blockedDesc, player.BlockedCharacters)
	} else {
		blocked = fmt.Sprintf("%s not yet set", blockedDesc)
	}
	backspaceText := "enabled"
	backspaceStyle := lipgloss.NewStyle().Foreground(styles.Green)
	if !player.CanDelete {
		backspaceText = "disabled"
		backspaceStyle = backspaceStyle.Foreground(styles.Red)
	}
	backspace := fmt.Sprintf("%s %s", backspaceDesc, backspaceStyle.Render(backspaceText))
	allowed = renderer.Render(allowed)
	blocked = renderer.Render(blocked)
	backspace = renderer.Render(backspace)
	return lipgloss.JoinVertical(lipgloss.Left, allowed, blocked, backspace)
}
