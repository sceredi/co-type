package styles

import (
	"charm.land/lipgloss/v2"
)

// Container is a style for a container with padding and a border.
var Container = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder(), true)

// NewContainer returns the given content styled with Container.
func NewContainer(content string) string {
	return Container.Render(content)
}
