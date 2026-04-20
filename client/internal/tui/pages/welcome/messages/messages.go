// Package welcome_messages defines the messages for the welcome page of the TUI.
package welcome_messages

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/common/domain"
)

// JoinLobbyMsg is a message that is sent when the user wants to join a lobby.
// It contains the ID of the lobby to join.
type JoinLobbyMsg struct {
	Lobby domain.Lobby
}

// JoinLobbyErrorMsg is a message that is sent when there is an error joining a lobby.
type JoinLobbyErrorMsg struct {
	Error string
}

// NewCreateLobbyCmd creates a new command to create a lobby.
func NewCreateLobbyCmd(lobbyCode string, playerName string) tea.Cmd {
	return func() tea.Msg {
		if lobbyCode == "" || playerName == "" {
			return JoinLobbyErrorMsg{
				Error: "Lobby code and username must not be empty",
			}
		}
		// TODO: actually create the lobby
		player := domain.NewPlayer(playerName)
		lobby := domain.NewLobby(lobbyCode, &player)
		return JoinLobbyMsg{Lobby: lobby}
	}
}

// NewJoinLobbyCmd creates a new command to join a lobby with the given ID.
func NewJoinLobbyCmd(lobbyCode string, playerName string) tea.Cmd {
	return func() tea.Msg {
		// TODO: try to join, then based on the response join or return an error
		if lobbyCode == "" || playerName == "" {
			return JoinLobbyErrorMsg{
				Error: "Lobby code and username must not be empty",
			}
		}
		log.Printf("Joining lobby %s as %s", lobbyCode, playerName)
		return JoinLobbyErrorMsg{Error: fmt.Sprintf("Unable to join lobby \"%s\"", lobbyCode)}
	}
}

// LeaveJoinLobbyMsg is a message that is sent when the user wants to leave the join lobby page and go back to the welcome page.
type LeaveJoinLobbyMsg struct{}

// NewLeaveJoinLobbyCmd creates a new command to leave the join lobby page and go back to the welcome page.
func NewLeaveJoinLobbyCmd() tea.Cmd {
	return func() tea.Msg { return LeaveJoinLobbyMsg{} }
}
