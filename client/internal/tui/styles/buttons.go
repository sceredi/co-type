// Package styles defines the styles for the TUI components.
package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	// ButtonDefault is a style for buttons with padding, border, and bold text.
	ButtonDefault = lipgloss.NewStyle().Padding(0, 3).Border(lipgloss.RoundedBorder()).Bold(true)
	// ButtonPrimary is a style for primary buttons with a green border and foreground.
	ButtonPrimary = ButtonDefault.BorderForeground(Green).Foreground(Green)
	// ButtonDanger is a style for danger buttons with a red border and foreground.
	ButtonDanger = ButtonDefault.BorderForeground(Red).Foreground(Red)
)
