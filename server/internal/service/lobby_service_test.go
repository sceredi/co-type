package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
	serverdomain "github.com/sceredi/co-type/server/internal/domain"
	"github.com/sceredi/co-type/server/internal/service"
)

type mockLobbyRepository struct {
	received *serverdomain.Lobby
	toReturn *serverdomain.Lobby
	getResp  *serverdomain.Lobby
	err      error
}

func (m *mockLobbyRepository) Create(lobby *serverdomain.Lobby) (*serverdomain.Lobby, error) {
	m.received = lobby
	if m.err != nil {
		return nil, m.err
	}
	if m.toReturn != nil {
		return m.toReturn, nil
	}
	return lobby, nil
}

func (m *mockLobbyRepository) Get(id string) *serverdomain.Lobby {
	return m.getResp
}

type mockControlGateway struct {
	receivedLobbyID    string
	receivedServerName string
	err                error
}

func (m *mockControlGateway) RegisterServer(name, host string, port int) error {
	return nil
}

func (m *mockControlGateway) RegisterLobby(lobbyID, serverName string) error {
	m.receivedLobbyID = lobbyID
	m.receivedServerName = serverName
	return m.err
}

func TestNewLobbyService(t *testing.T) {
	svc := service.NewLobbyService("test-server", &mockControlGateway{}, &mockLobbyRepository{})
	if svc == nil {
		t.Fatal("NewLobbyService() returned nil")
	}
}

func TestLobbyService_Create(t *testing.T) {
	repoErr := errors.New("repository create failed")
	gatewayErr := errors.New("gateway register failed")
	repoLobby := serverdomain.NewLobby("repo-lobby", playerPtr("repo-host"))
	serverName := "test-server"

	tests := []struct {
		name     string
		id       string
		userName string
		repo     *mockLobbyRepository
		gtw      *mockControlGateway
		wantErr  error
	}{
		{
			name:     "creates lobby and delegates to repository and gateway",
			id:       "lobby-1",
			userName: "alice",
			repo:     &mockLobbyRepository{toReturn: repoLobby},
			gtw:      &mockControlGateway{},
			wantErr:  nil,
		},
		{
			name:     "returns repository error",
			id:       "lobby-1",
			userName: "alice",
			repo:     &mockLobbyRepository{err: repoErr},
			gtw:      &mockControlGateway{},
			wantErr:  repoErr,
		},
		{
			name:     "returns gateway error",
			id:       "lobby-1",
			userName: "alice",
			repo:     &mockLobbyRepository{toReturn: repoLobby},
			gtw:      &mockControlGateway{err: gatewayErr},
			wantErr:  gatewayErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(serverName, tt.gtw, tt.repo)

			got, err := svc.Create(tt.id, tt.userName)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}

			if tt.repo.received == nil {
				t.Fatal("Create() repository was not called")
			}
			// ... existing checks ...
			if tt.repo.received.Base == nil {
				t.Fatal("Create() received nil Base, want non-nil base lobby")
			}
			if tt.repo.received.Base.ID != tt.id {
				t.Fatalf("Create() received id = %q, want %q", tt.repo.received.Base.ID, tt.id)
			}
			if tt.repo.received.Base.Host == nil || tt.repo.received.Base.Host.Name != tt.userName {
				t.Fatalf("Create() received host = %+v, want host name %q", tt.repo.received.Base.Host, tt.userName)
			}
			if len(tt.repo.received.Base.Players) != 1 || tt.repo.received.Base.Players[0] != tt.repo.received.Base.Host {
				t.Fatalf("Create() received players = %+v, want one host player", tt.repo.received.Base.Players)
			}

			if err == nil || errors.Is(err, gatewayErr) {
				// Search for repository calls only if repository should have been called
				if tt.gtw.receivedLobbyID != tt.id {
					t.Fatalf("Create() gateway received id = %q, want %q", tt.gtw.receivedLobbyID, tt.id)
				}
				if tt.gtw.receivedServerName != serverName {
					t.Fatalf("Create() gateway received serverName = %q, want %q", tt.gtw.receivedServerName, serverName)
				}
			}

			if err != nil {
				if got != nil {
					t.Fatalf("Create() got = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("Create() got nil, want non-nil lobby")
			}
			if tt.repo.toReturn == nil {
				t.Fatal("Create() repository returned nil lobby")
			}
			if tt.repo.toReturn.Subs == nil {
				t.Fatal("Create() repository returned lobby with nil Subs")
			}
			if got != tt.repo.toReturn {
				t.Fatalf("Create() got = %p, want %p", got, tt.repo.toReturn)
			}
			if got.Subs[tt.userName] == nil {
				t.Fatalf("Create() expected subscription channel for user %q", tt.userName)
			}
		})
	}
}

func TestLobbyService_Join(t *testing.T) {
	serverName := "test-server"
	lobbyID := "lobby-1"
	userName := "bob"
	existingLobby := serverdomain.NewLobby(lobbyID, playerPtr("alice"))

	tests := []struct {
		name    string
		repo    *mockLobbyRepository
		wantErr error
	}{
		{
			name:    "joins existing lobby successfully",
			repo:    &mockLobbyRepository{getResp: existingLobby},
			wantErr: nil,
		},
		{
			name:    "returns error if lobby not found",
			repo:    &mockLobbyRepository{getResp: nil},
			wantErr: domain.ErrLobbyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(serverName, &mockControlGateway{}, tt.repo)
			lobby, err := svc.Join(lobbyID, userName)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Join() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if lobby == nil {
					t.Fatal("expected lobby, got nil")
				}
				found := false
				for _, p := range lobby.Base.Players {
					if p.Name == userName {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected player %s in lobby players", userName)
				}
				if _, ok := lobby.Subs[userName]; !ok {
					t.Error("expected subscription channel for player")
				}
			}
		})
	}
}

func TestLobbyService_Get(t *testing.T) {
	want := serverdomain.NewLobby("lobby-1", playerPtr("alice"))
	repo := &mockLobbyRepository{getResp: want}
	svc := service.NewLobbyService("test-server", &mockControlGateway{}, repo)

	got := svc.Get("lobby-1")
	if got != want {
		t.Fatalf("Get() got = %p, want %p", got, want)
	}
}

func playerPtr(name string) *domain.Player {
	p := domain.NewPlayer(name)
	return p
}
