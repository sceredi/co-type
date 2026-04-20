// Package common contains reusable components and functions for the TUI pages.
package common

import "charm.land/bubbles/v2/textinput"

// NewTextinput creates a new text input with the given placeholder and value.
func NewTextinput(placeholder string, value string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.SetWidth(20)
	ti.SetValue(value)
	return ti
}
