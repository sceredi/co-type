package memory_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/broker/internal/domain"
	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/broker/internal/repository/memory"
)

func TestNewServerRepository(t *testing.T) {
	var repo repository.ServerRepository
	repo = memory.NewServerRepository()
	if repo == nil {
		t.Fatal("NewLobbyRepository() returned nil")
	}
}

func TestServerRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		seed    []*domain.Server
		input   *domain.Server
		wantErr error
	}{
		{
			name:    "creates new server",
			input:   &domain.Server{Addr: "8.8.8.8"},
			wantErr: nil,
		},
		{
			name: "returns ErrServerAlreadyExists when server with same name already exists",
			seed: []*domain.Server{
				{Name: "test_server"},
			},
			input:   &domain.Server{Name: "test_server"},
			wantErr: domain.ErrServerAlreadyExists,
		},
		{
			name: "returns ErrServerAlreadyExists when server with same addr already exists",
			seed: []*domain.Server{
				{Addr: "8.8.8.8"},
			},
			input:   &domain.Server{Addr: "8.8.8.8"},
			wantErr: domain.ErrServerAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewServerRepository()
			for _, server := range tt.seed {
				if _, err := repo.Create(server); err != nil {
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

			if got != tt.input {
				t.Fatalf("Create() got = %p, want same pointer as input = %p", got, tt.input)
			}
		})
	}
}
