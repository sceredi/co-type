// Package lobby defines the lobby page of the TUI.
package lobby

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/header"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/playerinfo"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/playerlist"
	"github.com/sceredi/co-type/client/internal/tui/pages/lobby/components/settings"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

type focusSlot int

const (
	focusPlayersList focusSlot = iota
	focusSettings
)

// The Model struct represents the state of the lobby page.
type Model struct {
	focus          focusSlot
	selectedPlayer int
	player         *domain.Player
	lobby          *domain.Lobby
	updates        <-chan *domain.Lobby

	ls service.LobbyService

	settings settings.Model
}

// New creates a new lobby model for the given player, lobby, and subscription channel.
func New(player *domain.Player, lobby *domain.Lobby, updates <-chan *domain.Lobby, ls service.LobbyService) Model {
	m := Model{
		focus:          focusPlayersList,
		selectedPlayer: 0,
		player:         player,
		lobby:          lobby,
		updates:        updates,
		ls:             ls,
	}
	return m
}

// Init initializes the model. It starts listening for lobby updates if a subscription is available.
func (m Model) Init() tea.Cmd {
	if m.updates == nil {
		return nil
	}
	return waitForLobbyUpdate(m.updates)
}

func waitForLobbyUpdate(updates <-chan *domain.Lobby) tea.Cmd {
	return func() tea.Msg {
		l, ok := <-updates
		if !ok {
			return lobby_messages.LobbySubscriptionClosedMsg{}
		}
		return lobby_messages.LobbyUpdatedMsg{Lobby: l}
	}
}

// Update updates the model based on the given message and returns the updated model and any commands to execute.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case lobby_messages.UpdatePlayerMsg:
		msg.Player.AllowedCharacters = msg.AllowedCharacters
		msg.Player.BlockedCharacters = msg.BlockedCharacters
		msg.Player.CanDelete = msg.BackspaceAllowed
	case lobby_messages.CloseSettingsMsg:
		m.focus = focusPlayersList
	case lobby_messages.LobbyUpdatedMsg:
		m.lobby = msg.Lobby
		cmds = append(cmds, waitForLobbyUpdate(m.updates))
	case lobby_messages.LobbySubscriptionClosedMsg:
		// subscription ended, no further updates
	}

	if m.focus == focusPlayersList {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "i":
				cmds = append(cmds, tea.SetClipboard(m.lobby.ID))
			case "up":
				m.handlePlayerSelect(m.selectedPlayer - 1)
			case "down":
				m.handlePlayerSelect(m.selectedPlayer + 1)
			case "r":
				cmd = m.handleReadyCmd()
			case "enter":
				m.handleOpenSettingsModal()
			case "esc":
				cmds = append(cmds, lobby_messages.NewLeaveLobbyCmd(m.ls, m.lobby.ID, m.player.Name))
			}
		}
	} else {
		m.settings, cmd = m.settings.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

const keybinds = "↑/↓: select player • enter: edit player\nr: toggle ready • i: copy lobby code\nesc: leave lobby • ctrl+c: quit"

// View renders the lobby page as a string.
func (m Model) View() string {
	readyLabel := "[ r ] Ready Up"
	readyButtonStyle := styles.ButtonDefault
	if m.player.IsReady {
		readyLabel = "[ r ] Ready!"
		readyButtonStyle = styles.ButtonPrimary
	}
	readyBtn := readyButtonStyle.Render(readyLabel)
	leaveBtn := styles.ButtonDanger.Render("[ esc ] Leave")

	body := lipgloss.JoinHorizontal(lipgloss.Center,
		playerlist.Render(m.selectedPlayer, m.lobby),
		playerinfo.Render(*m.lobby.Players[m.selectedPlayer]),
	)
	body = styles.NewContainer(body)

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, readyBtn, leaveBtn)

	v := lipgloss.JoinVertical(
		lipgloss.Center,
		header.Render(m.player, m.lobby),
		body,
		buttons,
		keybinds,
	)

	if m.focus == focusSettings {
		return m.settings.View()
	}

	return v
}
