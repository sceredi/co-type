package styles

import "charm.land/lipgloss/v2"

var (
	PaddingAll        = lipgloss.NewStyle().Padding(1)
	PaddingVertical   = lipgloss.NewStyle().Padding(1, 0)
	PaddingHorizontal = lipgloss.NewStyle().Padding(0, 1)
)
