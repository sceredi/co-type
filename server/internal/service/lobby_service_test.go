package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/service"
)

type mockLobbyRepository struct {
	received domain.Lobby
	toReturn *domain.Lobby
	err      error
}

func (m *mockLobbyRepository) Create(lobby domain.Lobby) (*domain.Lobby, error) {
	m.received = lobby
	if m.err != nil {
		return nil, m.err
	}
	if m.toReturn != nil {
		return m.toReturn, nil
	}
	return &lobby, nil
}

func TestNewLobbyService(t *testing.T) {
	svc := service.NewLobbyService(&mockLobbyRepository{})
	if svc == nil {
		t.Fatal("NewLobbyService() returned nil")
	}
}

func TestLobbyService_Create(t *testing.T) {
	repoErr := errors.New("repository create failed")
	repoLobby := domain.NewLobby("repo-lobby", playerPtr("repo-host"))

	tests := []struct {
		name     string
		id       string
		userName string
		repo     *mockLobbyRepository
		wantErr  error
	}{
		{
			name:     "creates lobby and delegates to repository",
			id:       "lobby-1",
			userName: "alice",
			repo:     &mockLobbyRepository{toReturn: &repoLobby},
			wantErr:  nil,
		},
		{
			name:     "returns repository error",
			id:       "lobby-1",
			userName: "alice",
			repo:     &mockLobbyRepository{err: repoErr},
			wantErr:  repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewLobbyService(tt.repo)

			got, err := svc.Create(tt.id, tt.userName)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}

			if tt.repo.received.ID != tt.id {
				t.Fatalf("Create() received id = %q, want %q", tt.repo.received.ID, tt.id)
			}
			if tt.repo.received.Host == nil || tt.repo.received.Host.Name != tt.userName {
				t.Fatalf("Create() received host = %+v, want host name %q", tt.repo.received.Host, tt.userName)
			}
			if len(tt.repo.received.Players) != 1 || tt.repo.received.Players[0] != tt.repo.received.Host {
				t.Fatalf("Create() received players = %+v, want one host player", tt.repo.received.Players)
			}

			if err != nil {
				if got != nil {
					t.Fatalf("Create() got = %v, want nil", got)
				}
				return
			}

			if got != tt.repo.toReturn {
				t.Fatalf("Create() got = %p, want %p", got, tt.repo.toReturn)
			}
		})
	}
}

func playerPtr(name string) *domain.Player {
	p := domain.NewPlayer(name)
	return &p
}
