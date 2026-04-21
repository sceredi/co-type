package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
	"github.com/sceredi/co-type/common/domain"
)

func TestJoinLobbyMsgSwitchesToLobbyPage(t *testing.T) {
	host := &domain.Player{Name: "Host"}
	lobby := domain.Lobby{
		ID:      "ABCD",
		Host:    host,
		Players: []*domain.Player{host},
	}

	m := New()
	updated, _ := m.Update(welcome_messages.JoinLobbyMsg{Lobby: lobby})
	got := updated.(Model)

	if got.page != lobbyPage {
		t.Fatalf("expected lobbyPage after join message, got %d", got.page)
	}
}

func TestLeaveMessageReturnsToWelcomePage(t *testing.T) {
	host := &domain.Player{Name: "Host"}
	lobby := domain.Lobby{
		ID:      "ABCD",
		Host:    host,
		Players: []*domain.Player{host},
	}

	m := New()
	updated, _ := m.Update(welcome_messages.JoinLobbyMsg{Lobby: lobby})
	m = updated.(Model)

	updated, _ = m.Update(lobby_messages.LeaveMsg{})
	got := updated.(Model)

	if got.page != welcomePage {
		t.Fatalf("expected welcomePage after leave message, got %d", got.page)
	}
}

func TestCtrlCReturnsQuitCommand(t *testing.T) {
	m := New()

	_, cmd := m.Update(tea.KeyPressMsg{Text: "ctrl+c"})
	if cmd == nil {
		t.Fatal("expected ctrl+c to return a quit command")
	}
}
