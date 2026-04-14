package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	buttonDefault = lipgloss.NewStyle().Padding(0, 3).Border(lipgloss.RoundedBorder()).Bold(true)
	ButtonPrimary = buttonDefault.BorderForeground(Green).Foreground(Green)
	ButtonDanger  = buttonDefault.BorderForeground(Red).Foreground(Red)
)
