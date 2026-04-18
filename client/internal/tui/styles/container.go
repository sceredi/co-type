package styles

import "charm.land/lipgloss/v2"

var Container = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder(), true)

func NewContainer(content string) string {
	return Container.Render(content)
}
