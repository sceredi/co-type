// Package game_messages defines the messages for the game page.
package game_messages

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/common/domain"
)

// GameEndMsg is a message that indicates that a game has ended.
type GameEndMsg struct {
	Stats domain.GameStats
}

// NewGameEndCmd returns a command that creates a GameEndMsg.
func NewGameEndCmd(stats domain.GameStats) tea.Cmd {
	return func() tea.Msg {
		return GameEndMsg{Stats: stats}
	}
}

// KeyPressResultMsg carries the lobby state returned after a key press, or an error.
type KeyPressResultMsg struct {
	Lobby *domain.Lobby
	Err   error
}

// NewSendKeyPressCmd returns a Cmd that sends a key press to the server and emits a KeyPressResultMsg.
func NewSendKeyPressCmd(lobbySvc service.LobbyService, lobbyID string, player *domain.Player, key string, isBackspace bool) tea.Cmd {
	return func() tea.Msg {
		l, err := lobbySvc.SendKeyPress(lobbyID, player, key, isBackspace)
		return KeyPressResultMsg{Lobby: l, Err: err}
	}
}

// RequestResumeGameResultMsg is emitted after the broker assigns a server for the crashed game.
type RequestResumeGameResultMsg struct {
	Server *domain.Server
	Err    error
}

// NewRequestResumeGameCmd contacts the broker to get a server that will host the resumed game.
func NewRequestResumeGameCmd(discoverySvc service.DiscoveryService, lobbyID string) tea.Cmd {
	return func() tea.Msg {
		srv, err := discoverySvc.RequestResumeGame(lobbyID)
		return RequestResumeGameResultMsg{Server: srv, Err: err}
	}
}

// ResumeGameResultMsg is emitted after the client has connected to the new server and resumed the game.
type ResumeGameResultMsg struct {
	Lobby   *domain.Lobby
	Updates <-chan *domain.Lobby
	Err     error
}

// NewResumeGameCmd connects to the given server, calls ResumeGame, and subscribes for updates.
func NewResumeGameCmd(lobbySvc service.LobbyService, server *domain.Server, lobby *domain.Lobby, playerName string) tea.Cmd {
	return func() tea.Msg {
		updatedLobby, updates, err := lobbySvc.ResumeGame(server, lobby, playerName)
		return ResumeGameResultMsg{Lobby: updatedLobby, Updates: updates, Err: err}
	}
}
