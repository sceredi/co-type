package lobby_messages

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/common/domain"
)

type PlayerReadyMsg struct {
	Name    string
	IsReady bool
}

func NewPlayerReadyCmd(name string, isReady bool) tea.Cmd {
	return func() tea.Msg {
		return PlayerReadyMsg{
			Name:    name,
			IsReady: isReady,
		}
	}
}

type UpdatePlayerMsg struct {
	Player            *domain.Player
	AllowedCharacters string
	BlockedCharacters string
	BackspaceAllowed  bool
}

func NewUpdatePlayerCmd(player *domain.Player, allowedCharacters, blockedCharacters string, backspaceAllowed bool) tea.Cmd {
	return func() tea.Msg {
		return UpdatePlayerMsg{
			Player:            player,
			AllowedCharacters: allowedCharacters,
			BlockedCharacters: blockedCharacters,
			BackspaceAllowed:  backspaceAllowed,
		}
	}
}

type CloseSettingsMsg struct{}

func NewCloseSettingsCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseSettingsMsg{}
	}
}
