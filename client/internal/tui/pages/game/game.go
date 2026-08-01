// Package game defines the game page of the TUI.
package game

import (
	"context"
	"log/slog"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/client/internal/tui/pages"
	game_messages "github.com/sceredi/co-type/client/internal/tui/pages/game/messages"
	"github.com/sceredi/co-type/client/internal/tui/styles"
	"github.com/sceredi/co-type/common/domain"
)

// Model is the model for the game page.
type Model struct {
	width  int
	height int

	lobby    domain.Lobby
	player   *domain.Player
	lobbySvc service.LobbyService
	discSvc  service.DiscoveryService
	updates  <-chan *domain.Lobby
}

// New creates a new game model.
func New(lobby domain.Lobby, player *domain.Player, updates <-chan *domain.Lobby, lobbySvc service.LobbyService, discSvc service.DiscoveryService) Model {
	return Model{lobby: lobby, player: player, updates: updates, lobbySvc: lobbySvc, discSvc: discSvc}
}

// Init initializes the game model.
func (m Model) Init() tea.Cmd {
	if m.updates == nil {
		return nil
	}
	return waitForGameUpdate(m.updates)
}

func waitForGameUpdate(updates <-chan *domain.Lobby) tea.Cmd {
	return func() tea.Msg {
		l, ok := <-updates
		if !ok {
			return serverDisconnectedMsg{}
		}
		return gameUpdateMsg{lobby: l}
	}
}

// serverDisconnectedMsg is emitted when the subscription channel is closed,
// indicating the server has become unreachable.
type serverDisconnectedMsg struct{}

type gameUpdateMsg struct{ lobby *domain.Lobby }

// Update updates the game model based on the given message.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		key := msg.String()
		isBackspace := key == "backspace"

		slog.DebugContext(context.Background(), "key pressed", slog.String("key", key))

		// Only send to server if the game is active.
		if m.lobby.Status == domain.LobbyPlaying {
			cmds = append(cmds, game_messages.NewSendKeyPressCmd(
				m.lobbySvc, m.lobby.ID, m.player, key, isBackspace,
			))
		}

	case game_messages.KeyPressResultMsg:
		if msg.Err != nil {
			slog.ErrorContext(context.Background(), "key press rejected", slog.String("err", msg.Err.Error()))
			break
		}
		m.lobby = *msg.Lobby
		for _, p := range m.lobby.Players {
			if p.Name == m.player.Name {
				m.player = p
				break
			}
		}
		if m.lobby.Status == domain.LobbyGameEnded {
			cmds = append(cmds, game_messages.NewGameEndCmd(domain.GameStats{
				TotalTime: time.Duration(m.lobby.Game.ElapsedMs) * time.Millisecond,
				Lobby:     m.lobby,
			}))
		}

	case serverDisconnectedMsg:
		slog.InfoContext(context.Background(), "server disconnected, pausing game and requesting resume")
		m.lobby.Status = domain.LobbyPaused
		cmds = append(cmds, game_messages.NewRequestResumeGameCmd(m.discSvc, m.lobby.ID))

	case game_messages.RequestResumeGameResultMsg:
		if msg.Err != nil {
			slog.ErrorContext(context.Background(), "failed to get resume server from broker", slog.String("err", msg.Err.Error()))
			break
		}
		slog.InfoContext(context.Background(), "broker assigned resume server", slog.String("server", msg.Server.Name))
		cmds = append(cmds, game_messages.NewResumeGameCmd(m.lobbySvc, msg.Server, &m.lobby, m.player.Name))

	case game_messages.ResumeGameResultMsg:
		if msg.Err != nil {
			slog.ErrorContext(context.Background(), "failed to resume game on new server", slog.String("err", msg.Err.Error()))
			break
		}
		slog.InfoContext(context.Background(), "game resumed on new server")
		m.lobby = *msg.Lobby
		m.updates = msg.Updates
		for _, p := range m.lobby.Players {
			if p.Name == m.player.Name {
				m.player = p
				break
			}
		}
		cmds = append(cmds, waitForGameUpdate(m.updates))

	case gameUpdateMsg:
		if msg.lobby != nil {
			m.lobby = *msg.lobby
			for _, p := range m.lobby.Players {
				if p.Name == m.player.Name {
					m.player = p
					break
				}
			}
			if m.lobby.Status == domain.LobbyGameEnded {
				cmds = append(cmds, game_messages.NewGameEndCmd(domain.GameStats{
					TotalTime: time.Duration(m.lobby.Game.ElapsedMs) * time.Millisecond,
					Lobby:     m.lobby,
				}))
			}
		}
		cmds = append(cmds, waitForGameUpdate(m.updates))
	}

	return m, tea.Batch(cmds...)
}

func (m Model) renderPlayers() string {
	rows := make([]string, len(m.lobby.Players)+1)
	rows[0] = "Players:"
	for i, p := range m.lobby.Players {
		rows[i+1] = p.Name
	}
	v := lipgloss.JoinVertical(lipgloss.Left, rows...)
	return v
}

func snippetPositions(snippet string, correctTypeable, wrongTypeable int64) (correctEnd, wrongEnd int) {
	var typed int64
	correctEnd = len(snippet)
	wrongEnd = len(snippet)
	totalTyped := correctTypeable + wrongTypeable

	for i, r := range snippet {
		if typed == correctTypeable && correctEnd == len(snippet) {
			correctEnd = i
		}
		if typed == totalTyped && wrongEnd == len(snippet) {
			wrongEnd = i
			break
		}
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			typed++
		}
	}
	return correctEnd, wrongEnd
}

// Custom rendering to avoid wrong new lines.
func renderInlineBackground(s lipgloss.Style, text string) string {
	lines := strings.Split(text, "\n")
	rendered := make([]string, len(lines))
	for i, line := range lines {
		rendered[i] = s.Render(line)
	}
	return strings.Join(rendered, "\n")
}

// View renders the game page.
func (m Model) View() string {
	style := lipgloss.NewStyle()
	snippet := m.lobby.Game.Snippet
	correctTypeable := m.lobby.Game.CorrectChars
	wrongTypeable := m.lobby.Game.WrongChars

	correctEnd, wrongEnd := snippetPositions(snippet, correctTypeable, wrongTypeable)

	correct := renderInlineBackground(style.Background(styles.DarkGreen), snippet[:correctEnd])
	wrong := renderInlineBackground(style.Background(styles.DarkRed), snippet[correctEnd:wrongEnd])
	rest := snippet[wrongEnd:]
	board := lipgloss.NewStyle().
		Width(max(m.width, pages.Width) * 4 / 5).
		Height(max(m.height, pages.Height) - 4).
		Align(lipgloss.Left).
		PaddingRight(2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		Render(correct + wrong + rest)
	players := lipgloss.NewStyle().PaddingLeft(2).Render(m.renderPlayers())

	v := lipgloss.JoinHorizontal(lipgloss.Top, board, players)

	v = styles.NewContainer(v)

	if m.lobby.Status == domain.LobbyPaused {
		base := lipgloss.NewLayer(v)
		topContent := styles.NewContainer(
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Game paused",
			),
		)
		x := (lipgloss.Width(v) / 2) - (lipgloss.Width(topContent) / 2)
		y := (lipgloss.Height(v) / 2) - (lipgloss.Height(topContent) / 2)
		top := lipgloss.NewLayer(
			topContent,
		).X(x).Y(y).Z(1)
		comp := lipgloss.NewCompositor(base, top)
		return comp.Render()
	}
	return v
}
