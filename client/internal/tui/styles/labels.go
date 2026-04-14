package styles

import "charm.land/lipgloss/v2"

var LabelBold = lipgloss.NewStyle().Bold(true)

func NewLabelBold(content string) string {
	return LabelBold.Render(content)
}
