package lobby

import tea "charm.land/bubbletea/v2"

type PlayerReadyMsg struct {
	Name    string
	IsReady bool
}

func NewPlayerReadyMsg(name string, isReady bool) tea.Cmd {
	return func() tea.Msg {
		return PlayerReadyMsg{
			Name:    name,
			IsReady: isReady,
		}
	}
}
