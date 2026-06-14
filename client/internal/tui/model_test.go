package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	lobby_messages "github.com/sceredi/co-type/client/internal/tui/pages/lobby/messages"
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

func (m *mockLobbyService) Connect(target *domain.Server) error {
	return nil
}

func newMockLobbyService() *mockLobbyService {
	return &mockLobbyService{}
}

func TestJoinLobbyMsgSwitchesToLobbyPage(t *testing.T) {
	host := &domain.Player{Name: "Host"}
	lobby := &domain.Lobby{
		ID:      "ABCD",
		Host:    host,
		Players: []*domain.Player{host},
	}
	m := New(newMockDiscoveryService(), newMockLobbyService())
	updated, _ := m.Update(welcome_messages.JoinedLobbyMsg{Lobby: lobby})
	got := updated.(Model)

	if got.page != lobbyPage {
		t.Fatalf("expected lobbyPage after join message, got %d", got.page)
	}
}

func TestLeaveMessageReturnsToWelcomePage(t *testing.T) {
	host := &domain.Player{Name: "Host"}
	lobby := &domain.Lobby{
		ID:      "ABCD",
		Host:    host,
		Players: []*domain.Player{host},
	}

	m := New(newMockDiscoveryService(), newMockLobbyService())
	updated, _ := m.Update(welcome_messages.JoinedLobbyMsg{Lobby: lobby})
	m = updated.(Model)

	updated, _ = m.Update(lobby_messages.LeaveMsg{})
	got := updated.(Model)

	if got.page != welcomePage {
		t.Fatalf("expected welcomePage after leave message, got %d", got.page)
	}
}

func TestCtrlCReturnsQuitCommand(t *testing.T) {
	m := New(newMockDiscoveryService(), newMockLobbyService())

	_, cmd := m.Update(tea.KeyPressMsg{Text: "ctrl+c"})
	if cmd == nil {
		t.Fatal("expected ctrl+c to return a quit command")
	}
}
