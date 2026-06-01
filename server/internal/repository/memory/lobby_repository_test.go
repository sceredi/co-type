package memory_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
	serverdomain "github.com/sceredi/co-type/server/internal/domain"
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
		seed    *serverdomain.Lobby
		input   *serverdomain.Lobby
		wantErr error
	}{
		{
			name:    "creates a new lobby",
			input:   serverdomain.NewLobby("lobby-1", domain.NewPlayer("alice")),
			wantErr: nil,
		},
		{
			name:    "returns ErrLobbyAlreadyExists when lobby id already exists",
			seed:    serverdomain.NewLobby("lobby-1", domain.NewPlayer("alice")),
			input:   serverdomain.NewLobby("lobby-1", domain.NewPlayer("bob")),
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
			if got.Base == nil {
				t.Fatal("Create() got nil Base, want non-nil base lobby")
			}
			if got.Base.ID != tt.input.Base.ID {
				t.Fatalf("Create() lobby id = %q, want %q", got.Base.ID, tt.input.Base.ID)
			}
			if got.Base.Host == nil || got.Base.Host.Name != tt.input.Base.Host.Name {
				t.Fatalf("Create() host = %+v, want host name %q", got.Base.Host, tt.input.Base.Host.Name)
			}
			if len(got.Base.Players) != 1 || got.Base.Players[0] != got.Base.Host {
				t.Fatalf("Create() players = %+v, want one host player", got.Base.Players)
			}
		})
	}
}
