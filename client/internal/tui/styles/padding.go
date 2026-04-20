package styles

import "charm.land/lipgloss/v2"

var (
	// PaddingAll is a style with padding of 1 on all sides.
	PaddingAll = lipgloss.NewStyle().Padding(1)
	// PaddingVertical is a style with padding of 1 on the top and bottom.
	PaddingVertical = lipgloss.NewStyle().Padding(1, 0)
	// PaddingHorizontal is a style with padding of 1 on the left and right.
	PaddingHorizontal = lipgloss.NewStyle().Padding(0, 1)
)
