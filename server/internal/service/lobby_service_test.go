package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
	serverdomain "github.com/sceredi/co-type/server/internal/domain"
	"github.com/sceredi/co-type/server/internal/service"
)

type mockLobbyRepository struct {
	received  *serverdomain.Lobby
	toReturn  *serverdomain.Lobby
	getResp   *serverdomain.Lobby
	err       error
	deleteErr error
	deletedID string
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

func (m *mockLobbyRepository) Delete(id string) error {
	m.deletedID = id
	return m.deleteErr
}

type mockControlGateway struct {
	receivedLobbyID    string
	receivedServerName string
	unregisteredID     string
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

func (m *mockControlGateway) UnregisterLobby(lobbyID string) error {
	m.unregisteredID = lobbyID
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
			if len(tt.repo.received.Base.Players) != 1 || tt.repo.received.Base.Players[0].Name != tt.userName {
				t.Fatalf("Create() received players = %+v, want one player with name %q", tt.repo.received.Base.Players, tt.userName)
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

func TestLobbyService_Create_DeleteOnGatewayError(t *testing.T) {
	gatewayErr := errors.New("gateway register failed")
	deleteErr := errors.New("delete failed")
	repoLobby := serverdomain.NewLobby("lobby-1", playerPtr("alice"))
	serverName := "test-server"
	lobbyID := "lobby-1"

	tests := []struct {
		name          string
		repo          *mockLobbyRepository
		gtw           *mockControlGateway
		wantDeletedID string
	}{
		{
			name:          "calls Delete when gateway registration fails",
			repo:          &mockLobbyRepository{toReturn: repoLobby},
			gtw:           &mockControlGateway{err: gatewayErr},
			wantDeletedID: lobbyID,
		},
		{
			name:          "calls Delete when gateway fails even if Delete itself errors",
			repo:          &mockLobbyRepository{toReturn: repoLobby, deleteErr: deleteErr},
			gtw:           &mockControlGateway{err: gatewayErr},
			wantDeletedID: lobbyID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(serverName, tt.gtw, tt.repo)

			got, err := svc.Create(lobbyID, "alice")
			if !errors.Is(err, gatewayErr) {
				t.Fatalf("Create() error = %v, want %v", err, gatewayErr)
			}
			if got != nil {
				t.Fatalf("Create() got = %v, want nil", got)
			}
			if tt.repo.deletedID != tt.wantDeletedID {
				t.Fatalf("Create() deleted id = %q, want %q", tt.repo.deletedID, tt.wantDeletedID)
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
		{
			name:    "returns error if player name already in lobby",
			repo:    &mockLobbyRepository{getResp: serverdomain.NewLobby(lobbyID, playerPtr(userName))},
			wantErr: domain.ErrPlayerAlreadyInLobby,
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

func TestLobbyService_Leave(t *testing.T) {
	serverName := "test-server"
	lobbyID := "lobby-1"

	tests := []struct {
		name           string
		userName       string
		repo           *mockLobbyRepository
		gtw            *mockControlGateway
		wantErr        error
		wantUnregister bool
		wantDeletedID  string
	}{
		{
			name:     "returns error if lobby not found",
			userName: "alice",
			repo:     &mockLobbyRepository{getResp: nil},
			gtw:      &mockControlGateway{},
			wantErr:  domain.ErrLobbyNotFound,
		},
		{
			name:     "returns error if player not in lobby",
			userName: "charlie",
			repo:     &mockLobbyRepository{getResp: serverdomain.NewLobby(lobbyID, playerPtr("alice"))},
			gtw:      &mockControlGateway{},
			wantErr:  domain.ErrPlayerNotInLobby,
		},
		{
			name:           "removes player and notifies remaining players",
			userName:       "bob",
			repo:           &mockLobbyRepository{getResp: buildLobbyWithTwo(lobbyID, "alice", "bob")},
			gtw:            &mockControlGateway{},
			wantErr:        nil,
			wantUnregister: false,
		},
		{
			name:           "deletes and unregisters lobby when last player leaves",
			userName:       "alice",
			repo:           &mockLobbyRepository{getResp: serverdomain.NewLobby(lobbyID, playerPtr("alice"))},
			gtw:            &mockControlGateway{},
			wantErr:        nil,
			wantUnregister: true,
			wantDeletedID:  lobbyID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(serverName, tt.gtw, tt.repo)
			err := svc.Leave(lobbyID, tt.userName)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Leave() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantUnregister && tt.gtw.unregisteredID != lobbyID {
				t.Fatalf("Leave() expected UnregisterLobby to be called with %q, got %q", lobbyID, tt.gtw.unregisteredID)
			}
			if !tt.wantUnregister && tt.gtw.unregisteredID != "" {
				t.Fatalf("Leave() expected UnregisterLobby not to be called, but got %q", tt.gtw.unregisteredID)
			}
			if tt.wantDeletedID != "" && tt.repo.deletedID != tt.wantDeletedID {
				t.Fatalf("Leave() expected repo.Delete(%q), got %q", tt.wantDeletedID, tt.repo.deletedID)
			}
		})
	}
}

func buildLobbyWithTwo(id, host, other string) *serverdomain.Lobby {
	l := serverdomain.NewLobby(id, playerPtr(host))
	otherPlayer := domain.NewPlayer(other)
	l.Base.AddPlayers(otherPlayer)
	l.Subs[other] = make(chan *serverdomain.Lobby, 64)
	return l
}

func TestLobbyService_EditPlayer(t *testing.T) {
	serverName := "test-server"
	lobbyID := "lobby-1"

	tests := []struct {
		name             string
		playerName       string
		repo             *mockLobbyRepository
		wantErr          error
		wantIsReady      bool
		wantAllowedChars string
		wantBlockedChars string
		wantCanDelete    bool
	}{
		{
			name:       "returns error if lobby not found",
			playerName: "alice",
			repo:       &mockLobbyRepository{getResp: nil},
			wantErr:    domain.ErrLobbyNotFound,
		},
		{
			name:       "returns error if player not in lobby",
			playerName: "charlie",
			repo:       &mockLobbyRepository{getResp: serverdomain.NewLobby(lobbyID, playerPtr("alice"))},
			wantErr:    domain.ErrPlayerNotInLobby,
		},
		{
			name:             "updates player settings and notifies subscribers",
			playerName:       "alice",
			repo:             &mockLobbyRepository{getResp: buildLobbyWithSub(lobbyID, "alice")},
			wantErr:          nil,
			wantIsReady:      true,
			wantAllowedChars: "abc",
			wantBlockedChars: "xyz",
			wantCanDelete:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(serverName, &mockControlGateway{}, tt.repo)
			l, err := svc.EditPlayer(lobbyID, tt.playerName, tt.wantIsReady, tt.wantAllowedChars, tt.wantBlockedChars, tt.wantCanDelete)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("EditPlayer() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if l != nil {
					t.Fatalf("EditPlayer() got = %v, want nil", l)
				}
				return
			}

			if l == nil {
				t.Fatal("EditPlayer() got nil, want non-nil lobby")
			}

			var found *domain.Player
			for _, p := range l.Base.Players {
				if p.Name == tt.playerName {
					found = p
					break
				}
			}
			if found == nil {
				t.Fatalf("EditPlayer() player %q not found in lobby", tt.playerName)
			}
			if found.IsReady != tt.wantIsReady {
				t.Errorf("EditPlayer() IsReady = %v, want %v", found.IsReady, tt.wantIsReady)
			}
			if found.AllowedCharacters != tt.wantAllowedChars {
				t.Errorf("EditPlayer() AllowedCharacters = %q, want %q", found.AllowedCharacters, tt.wantAllowedChars)
			}
			if found.BlockedCharacters != tt.wantBlockedChars {
				t.Errorf("EditPlayer() BlockedCharacters = %q, want %q", found.BlockedCharacters, tt.wantBlockedChars)
			}
			if found.CanDelete != tt.wantCanDelete {
				t.Errorf("EditPlayer() CanDelete = %v, want %v", found.CanDelete, tt.wantCanDelete)
			}
		})
	}
}

func buildLobbyWithSub(id, playerName string) *serverdomain.Lobby {
	l := serverdomain.NewLobby(id, playerPtr(playerName))
	l.Subs[playerName] = make(chan *serverdomain.Lobby, 64)
	return l
}

func TestLobbyService_Ready(t *testing.T) {
	serverName := "test-server"
	lobbyID := "lobby-1"

	tests := []struct {
		name        string
		playerName  string
		repo        *mockLobbyRepository
		wantErr     error
		wantIsReady bool
		wantSnippet bool
	}{
		{
			name:       "returns error if lobby not found",
			playerName: "alice",
			repo:       &mockLobbyRepository{getResp: nil},
			wantErr:    domain.ErrLobbyNotFound,
		},
		{
			name:       "returns error if player not in lobby",
			playerName: "charlie",
			repo:       &mockLobbyRepository{getResp: buildLobbyWithSub(lobbyID, "alice")},
			wantErr:    domain.ErrPlayerNotInLobby,
		},
		{
			name:        "toggles player to ready and notifies subscribers",
			playerName:  "alice",
			repo:        &mockLobbyRepository{getResp: buildLobbyWithSub(lobbyID, "alice")},
			wantErr:     nil,
			wantIsReady: true,
			wantSnippet: true,
		},
		{
			name:        "sets snippet when all players are ready",
			playerName:  "bob",
			repo:        &mockLobbyRepository{getResp: buildReadyLobbyMissingOne(lobbyID, "alice", "bob")},
			wantErr:     nil,
			wantIsReady: true,
			wantSnippet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(serverName, &mockControlGateway{}, tt.repo)
			l, err := svc.Ready(lobbyID, tt.playerName)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Ready() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if l != nil {
					t.Fatalf("Ready() got = %v, want nil", l)
				}
				return
			}

			if l == nil {
				t.Fatal("Ready() got nil, want non-nil lobby")
			}

			var found *domain.Player
			for _, p := range l.Base.Players {
				if p.Name == tt.playerName {
					found = p
					break
				}
			}
			if found == nil {
				t.Fatalf("Ready() player %q not found in lobby", tt.playerName)
			}
			if found.IsReady != tt.wantIsReady {
				t.Errorf("Ready() IsReady = %v, want %v", found.IsReady, tt.wantIsReady)
			}
			if tt.wantSnippet && l.Base.Snippet == "" {
				t.Error("Ready() expected snippet to be set when all players are ready")
			}
		})
	}
}

// buildReadyLobbyMissingOne builds a lobby where all players except one are already ready.
func buildReadyLobbyMissingOne(id, readyPlayer, notReadyPlayer string) *serverdomain.Lobby {
	l := serverdomain.NewLobby(id, playerPtr(readyPlayer))
	l.Base.Players[0].IsReady = true
	notReady := domain.NewPlayer(notReadyPlayer)
	l.Base.AddPlayers(notReady)
	l.Subs[readyPlayer] = make(chan *serverdomain.Lobby, 64)
	l.Subs[notReadyPlayer] = make(chan *serverdomain.Lobby, 64)
	return l
}
