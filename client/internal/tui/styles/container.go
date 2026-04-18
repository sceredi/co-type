package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var Container = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder(), true)

func NewContainer(content string, color color.Color) string {
	return Container.BorderForeground(color).Render(content)
}
