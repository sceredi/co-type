package styles

import "charm.land/lipgloss/v2"

// LabelBold is a style for labels that are bold.
var LabelBold = lipgloss.NewStyle().Bold(true)

// NewLabelBold returns the given content styled with LabelBold.
func NewLabelBold(content string) string {
	return LabelBold.Render(content)
}
