package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	ButtonDefault = lipgloss.NewStyle().Padding(0, 3).Border(lipgloss.RoundedBorder()).Bold(true)
	ButtonPrimary = ButtonDefault.BorderForeground(Green).Foreground(Green)
	ButtonDanger  = ButtonDefault.BorderForeground(Red).Foreground(Red)
)
