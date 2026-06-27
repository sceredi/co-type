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
			if len(got.Base.Players) != 1 || got.Base.Players[0].Name != tt.input.Base.Players[0].Name {
				t.Fatalf("Create() players = %+v, want one player matching input", got.Base.Players)
			}
		})
	}
}

func TestLobbyRepository_Get(t *testing.T) {
	tests := []struct {
		name    string
		seed    *serverdomain.Lobby
		queryID string
		wantNil bool
	}{
		{
			name:    "returns existing lobby",
			seed:    serverdomain.NewLobby("lobby-1", domain.NewPlayer("alice")),
			queryID: "lobby-1",
			wantNil: false,
		},
		{
			name:    "returns nil for missing lobby",
			queryID: "missing",
			wantNil: true,
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

			got := repo.Get(tt.queryID)
			if tt.wantNil && got != nil {
				t.Fatalf("Get() got = %v, want nil", got)
			}
			if !tt.wantNil && got == nil {
				t.Fatal("Get() got nil, want non-nil lobby")
			}
		})
	}
}

func TestLobbyRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		seed    *serverdomain.Lobby
		queryID string
		wantErr error
	}{
		{
			name:    "deletes existing lobby",
			seed:    serverdomain.NewLobby("lobby-1", domain.NewPlayer("alice")),
			queryID: "lobby-1",
			wantErr: nil,
		},
		{
			name:    "returns ErrLobbyNotFound for missing lobby",
			queryID: "missing",
			wantErr: domain.ErrLobbyNotFound,
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

			err := repo.Delete(tt.queryID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if got := repo.Get(tt.queryID); got != nil {
					t.Fatalf("Delete() lobby still exists after deletion")
				}
			}
		})
	}
}

func TestLobbyRepository_Delete_Twice(t *testing.T) {
	repo := memory.NewLobbyRepository()
	l := serverdomain.NewLobby("lobby-1", domain.NewPlayer("alice"))
	if _, err := repo.Create(l); err != nil {
		t.Fatalf("seed create failed: %v", err)
	}

	if err := repo.Delete("lobby-1"); err != nil {
		t.Fatalf("first Delete() unexpected error: %v", err)
	}
	if err := repo.Delete("lobby-1"); !errors.Is(err, domain.ErrLobbyNotFound) {
		t.Fatalf("second Delete() error = %v, want ErrLobbyNotFound", err)
	}
}

func TestLobbyRepository_Delete_DoesNotAffectOtherLobbies(t *testing.T) {
	repo := memory.NewLobbyRepository()
	l1 := serverdomain.NewLobby("lobby-1", domain.NewPlayer("alice"))
	l2 := serverdomain.NewLobby("lobby-2", domain.NewPlayer("bob"))
	if _, err := repo.Create(l1); err != nil {
		t.Fatalf("create lobby-1 failed: %v", err)
	}
	if _, err := repo.Create(l2); err != nil {
		t.Fatalf("create lobby-2 failed: %v", err)
	}

	if err := repo.Delete("lobby-1"); err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}
	if got := repo.Get("lobby-2"); got == nil {
		t.Fatal("Delete() removed lobby-2, want it to still exist")
	}
}
