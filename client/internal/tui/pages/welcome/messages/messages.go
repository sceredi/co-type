// Package welcome_messages defines the messages for the welcome page of the TUI.
package welcome_messages

import (
	"context"
	"fmt"
	"log/slog"

	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/common/domain"
)

// JoinedLobbyMsg is a message that is sent when the user wants to join a lobby.
type JoinedLobbyMsg struct {
	Lobby *domain.Lobby
}

// JoinLobbyErrorMsg is a message that is sent when there is an error joining a lobby.
type JoinLobbyErrorMsg struct {
	Error string
}

// NewCreateLobbyCmd creates a new command to create a lobby.
func NewCreateLobbyCmd(ds service.DiscoveryService, ls service.LobbyService, lobbyCode string, playerName string) tea.Cmd {
	return func() tea.Msg {
		if lobbyCode == "" || playerName == "" {
			return JoinLobbyErrorMsg{
				Error: "Lobby code and username must not be empty",
			}
		}
		server, err := ds.GetAvailableServer()
		slog.InfoContext(context.Background(), "got server",
			slog.Any("server", server),
		)
		if err != nil {
			slog.ErrorContext(context.Background(), "unable to get server",
				slog.String("err", err.Error()),
			)
			return JoinLobbyErrorMsg{
				Error: fmt.Sprintf("Unable to get a game server:\n%s", err.Error()),
			}
		}
		err = ls.Connect(server)
		if err != nil {
			slog.ErrorContext(context.Background(), "unable to connect to gameserver",
				slog.String("err", err.Error()),
			)
			return JoinLobbyErrorMsg{
				Error: fmt.Sprintf("Unable to connect to game server:\n%s", err.Error()),
			}
		}
		lobby, err := ls.Create(lobbyCode, playerName)
		if err != nil {
			slog.ErrorContext(context.Background(), "unable to create lobby",
				slog.String("err", err.Error()),
			)
			return JoinLobbyErrorMsg{
				Error: fmt.Sprintf("Unable to create lobby:\n%s", err.Error()),
			}
		}
		return JoinedLobbyMsg{Lobby: lobby}
	}
}

// NewJoinLobbyCmd creates a new command to join a lobby with the given ID.
func NewJoinLobbyCmd(ds service.DiscoveryService, ls service.LobbyService, lobbyCode string, playerName string) tea.Cmd {
	return func() tea.Msg {
		// TODO: try to join, then based on the response join or return an error
		if lobbyCode == "" || playerName == "" {
			return JoinLobbyErrorMsg{
				Error: "Lobby code and username must not be empty",
			}
		}
		server, err := ds.GetHostByLobby(lobbyCode)
		slog.InfoContext(context.Background(), "got server",
			slog.Any("server", server),
		)
		if err != nil {
			slog.ErrorContext(context.Background(), "unable to get server",
				slog.String("err", err.Error()),
			)
			return JoinLobbyErrorMsg{
				Error: fmt.Sprintf("Unable to get a game server:\n%s", err.Error()),
			}
		}
		err = ls.Connect(server)
		if err != nil {
			slog.ErrorContext(context.Background(), "unable to connect to gameserver",
				slog.String("err", err.Error()),
			)
			return JoinLobbyErrorMsg{
				Error: fmt.Sprintf("Unable to connect to game server:\n%s", err.Error()),
			}
		}
		lobby, err := ls.Join(lobbyCode, playerName)
		if err != nil {
			slog.ErrorContext(context.Background(), "unable to join lobby",
				slog.String("err", err.Error()),
			)
			return JoinLobbyErrorMsg{
				Error: fmt.Sprintf("Unable to join lobby:\n%s", err.Error()),
			}
		}
		return JoinedLobbyMsg{Lobby: lobby}
	}
}

// LeaveJoinLobbyMsg is a message that is sent when the user wants to leave the join lobby page and go back to the welcome page.
type LeaveJoinLobbyMsg struct{}

// NewLeaveJoinLobbyCmd creates a new command to leave the join lobby page and go back to the welcome page.
func NewLeaveJoinLobbyCmd() tea.Cmd {
	return func() tea.Msg { return LeaveJoinLobbyMsg{} }
}
