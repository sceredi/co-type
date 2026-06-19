package welcome

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	welcome_messages "github.com/sceredi/co-type/client/internal/tui/pages/welcome/messages"
	"github.com/sceredi/co-type/common/domain"
)

// mockDiscoveryService is a mock implementation of DiscoveryService for testing
type mockDiscoveryService struct{}

func (m *mockDiscoveryService) GetAvailableServer() (*domain.Server, error) {
	return nil, nil
}

func (m *mockDiscoveryService) GetHostByLobby(lobbyCode string) (*domain.Server, error) {
	return nil, nil
}

func newMockDiscoveryService() *mockDiscoveryService {
	return &mockDiscoveryService{}
}

// mockLobbyService is a mock implementation of LobbyService for testing
type mockLobbyService struct{}

func (m *mockLobbyService) Create(id, hostName string) (*domain.Lobby, error) {
	host := &domain.Player{Name: hostName}
	return &domain.Lobby{ID: id, Host: host, Players: []*domain.Player{host}}, nil
}

func (m *mockLobbyService) Join(id, playerName string) (*domain.Lobby, error) {
	player := &domain.Player{Name: playerName}
	return &domain.Lobby{ID: id, Players: []*domain.Player{player}}, nil
}

func (m *mockLobbyService) Leave(id, playerName string) error {
	return nil
}

func (m *mockLobbyService) Connect(target *domain.Server) error {
	return nil
}

func (m *mockLobbyService) Subscribe(lobbyID, playerName string) (<-chan *domain.Lobby, error) {
	return nil, nil
}

func newMockLobbyService() *mockLobbyService {
	return &mockLobbyService{}
}

func TestTabCyclesThroughFocusSlots(t *testing.T) {
	m := New(newMockDiscoveryService(), newMockLobbyService())

	if m.focus != focusName {
		t.Fatalf("expected initial focusName, got %d", m.focus)
	}
	if !m.nameTi.Focused() {
		t.Fatal("expected name input to start focused")
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusCode {
		t.Fatalf("expected focusCode after first tab, got %d", m.focus)
	}
	if !m.codeTi.Focused() {
		t.Fatal("expected code input to be focused after first tab")
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusCreate {
		t.Fatalf("expected focusCreate after second tab, got %d", m.focus)
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusJoin {
		t.Fatalf("expected focusJoin after third tab, got %d", m.focus)
	}

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusName {
		t.Fatalf("expected focusName after wrapping, got %d", m.focus)
	}
	if !m.nameTi.Focused() {
		t.Fatal("expected name input to be focused after wrapping")
	}
}

func TestEnterOnNameAdvancesToCode(t *testing.T) {
	m := New(newMockDiscoveryService(), newMockLobbyService())

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if m.focus != focusCode {
		t.Fatalf("expected focusCode after enter on name, got %d", m.focus)
	}
}

func TestJoinLobbyErrorMsgUpdatesViewError(t *testing.T) {
	m := New(newMockDiscoveryService(), newMockLobbyService())
	m, _ = m.Update(welcome_messages.JoinLobbyErrorMsg{Error: "Unable to join lobby \"ABCD\""})

	if m.error != "Unable to join lobby \"ABCD\"" {
		t.Fatalf("expected error message to be stored, got %q", m.error)
	}
	if !strings.Contains(m.View(), m.error) {
		t.Fatalf("expected rendered view to include error %q", m.error)
	}
}
