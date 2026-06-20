package lobby

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
	"github.com/sceredi/co-type/common/domain"
)

func newTestLobbyModel() Model {
	host := &domain.Player{Name: "Host"}
	guest := &domain.Player{Name: "Guest"}
	third := &domain.Player{Name: "Third"}

	l := &domain.Lobby{
		ID:      "ABCD",
		Players: []*domain.Player{host, guest, third},
	}

	return New(host, l, nil, nil)
}

func TestDownMovesSelectionForward(t *testing.T) {
	m := newTestLobbyModel()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.selectedPlayer != 1 {
		t.Errorf("expected selectedPlayer=1, got %d", m.selectedPlayer)
	}
}

func TestUpWrapsSelectionToLastPlayer(t *testing.T) {
	m := newTestLobbyModel()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.selectedPlayer != 2 {
		t.Errorf("expected selectedPlayer=2, got %d", m.selectedPlayer)
	}
}

func TestDownWrapsSelectionToFirstPlayer(t *testing.T) {
	m := newTestLobbyModel()
	m.selectedPlayer = 2

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.selectedPlayer != 0 {
		t.Errorf("expected selectedPlayer=0, got %d", m.selectedPlayer)
	}
}

func TestEnterOpensSettingsModal(t *testing.T) {
	m := newTestLobbyModel()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.focus != focusSettings {
		t.Fatalf("expected focusSettings, got %d", m.focus)
	}
}

func TestCloseSettingsReturnsPlayersFocus(t *testing.T) {
	m := newTestLobbyModel()
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.focus != focusSettings {
		t.Fatalf("expected focusSettings before closing, got %d", m.focus)
	}

	m, _ = m.Update(lobby_messages.CloseSettingsMsg{})
	if m.focus != focusPlayersList {
		t.Fatalf("expected focusPlayersList after close, got %d", m.focus)
	}
}

func TestRTogglesCurrentPlayerReady(t *testing.T) {
	m := newTestLobbyModel()
	if m.player.IsReady {
		t.Fatal("expected player to start not ready")
	}

	// Pressing "r" dispatches a command to the server; the state update arrives
	// via LobbyUpdatedMsg. Simulate the server response by sending an updated lobby.
	m, _ = m.Update(tea.KeyPressMsg{Text: "r"})

	updatedLobby := &domain.Lobby{
		ID:      m.lobby.ID,
		Players: []*domain.Player{{Name: "Host", IsReady: true}, {Name: "Guest"}, {Name: "Third"}},
	}
	m, _ = m.Update(lobby_messages.LobbyUpdatedMsg{Lobby: updatedLobby})
	if !m.player.IsReady {
		t.Error("expected player to be ready after lobby update")
	}

	m, _ = m.Update(tea.KeyPressMsg{Text: "r"})

	updatedLobby2 := &domain.Lobby{
		ID:      m.lobby.ID,
		Players: []*domain.Player{{Name: "Host", IsReady: false}, {Name: "Guest"}, {Name: "Third"}},
	}
	m, _ = m.Update(lobby_messages.LobbyUpdatedMsg{Lobby: updatedLobby2})
	if m.player.IsReady {
		t.Error("expected player to be not ready after second lobby update")
	}
}
