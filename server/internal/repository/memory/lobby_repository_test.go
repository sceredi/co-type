package memory_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/server/internal/repository"
	"github.com/sceredi/co-type/server/internal/repository/memory"
)

func TestNewLobbyRepository(t *testing.T) {
	var repo repository.LobbyRepository
	repo = memory.NewLobbyRepository()
	if repo == nil {
		t.Fatal("NewLobbyRepository() returned nil")
	}
}

func TestLobbyRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		seed    *domain.Lobby
		input   *domain.Lobby
		wantErr error
	}{
		{
			name:    "creates a new lobby",
			input:   domain.NewLobby("lobby-1", domain.NewPlayer("alice")),
			wantErr: nil,
		},
		{
			name:    "returns ErrLobbyAlreadyExists when lobby id already exists",
			seed:    domain.NewLobby("lobby-1", domain.NewPlayer("alice")),
			input:   domain.NewLobby("lobby-1", domain.NewPlayer("bob")),
			wantErr: domain.ErrLobbyAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewLobbyRepository()

			if tt.seed != nil {
				if _, err := repo.Create(tt.seed); err != nil {
					t.Fatalf("seed create failed: %v", err)
				}
			}

			got, err := repo.Create(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
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
			if got.ID != tt.input.ID {
				t.Fatalf("Create() lobby id = %q, want %q", got.ID, tt.input.ID)
			}
			if got.Host == nil || got.Host.Name != tt.input.Host.Name {
				t.Fatalf("Create() host = %+v, want host name %q", got.Host, tt.input.Host.Name)
			}
			if len(got.Players) != 1 || got.Players[0] != got.Host {
				t.Fatalf("Create() players = %+v, want one host player", got.Players)
			}
		})
	}
}
