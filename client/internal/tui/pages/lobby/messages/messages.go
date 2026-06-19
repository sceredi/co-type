// Package lobby_messages defines the messages for the lobby page.
package lobby_messages

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/common/domain"
)

// PlayerReadyMsg is a message that indicates that a player has changed their ready status.
type PlayerReadyMsg struct {
	Name    string
	IsReady bool
}

// NewPlayerReadyCmd returns a command that creates a PlayerReadyMsg with the given name and ready status.
func NewPlayerReadyCmd(name string, isReady bool) tea.Cmd {
	return func() tea.Msg {
		return PlayerReadyMsg{
			Name:    name,
			IsReady: isReady,
		}
	}
}

// LeaveMsg is a message that indicates the player wants to leave the lobby and return to the main menu.
type LeaveMsg struct{}

// NewLeaveCmd returns a command that creates a LeaveMessage.
func NewLeaveCmd() tea.Cmd {
	return func() tea.Msg { return LeaveMsg{} }
}

// UpdatePlayerMsg is a message that indicates that a player's settings have been updated.
type UpdatePlayerMsg struct {
	Player            *domain.Player
	AllowedCharacters string
	BlockedCharacters string
	BackspaceAllowed  bool
}

// NewUpdatePlayerCmd returns a command that creates an UpdatePlayerMsg with the given player and settings.
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

// CloseSettingsMsg is a message that indicates that the settings modal should be closed.
type CloseSettingsMsg struct{}

// NewCloseSettingsCmd returns a command that creates a CloseSettingsMsg.
func NewCloseSettingsCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseSettingsMsg{}
	}
}

// LobbyUpdatedMsg is a message that carries an updated lobby state from the subscription stream.
type LobbyUpdatedMsg struct {
	Lobby *domain.Lobby
}

// LobbySubscriptionClosedMsg is a message that indicates the lobby subscription stream has closed.
type LobbySubscriptionClosedMsg struct{}

// StartGameMsg is a message that indicates that a game should be started.
type StartGameMsg struct {
	Game domain.Game
}

// NewStartGameCmd return a command that creates a StartGameMsg.
func NewStartGameCmd(game domain.Game) tea.Cmd {
	return func() tea.Msg {
		return StartGameMsg{Game: game}
	}
}
