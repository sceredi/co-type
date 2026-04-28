package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/broker/internal/domain"
	"github.com/sceredi/co-type/broker/internal/repository/memory"
	"github.com/sceredi/co-type/broker/internal/service"
)

func TestNewServerService(t *testing.T) {
	var svc service.ServerService
	repo := memory.NewServerRepository()
	svc = service.NewServerService(repo)
	if svc == nil {
		t.Fatal("NewServerService() returned nil")
	}
}

func TestServerService_CreateServer(t *testing.T) {
	tests := []struct {
		name    string
		seed    []*domain.Server
		addr    string
		port    int32
		wantErr error
	}{
		{
			name:    "creates new server",
			addr:    "8.8.8.8",
			port:    8080,
			wantErr: nil,
		},
		{
			name: "returns ErrServerAlreadyExists when server already exists",
			seed: []*domain.Server{
				{Addr: "8.8.8.8", Port: 8080},
			},
			addr:    "8.8.8.8",
			port:    9090,
			wantErr: domain.ErrServerAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.NewServerRepository()
			svc := service.NewServerService(repo)
			for _, server := range tt.seed {
				if _, err := repo.Create(server); err != nil {
					t.Fatalf("seed create failed: %v", err)
				}
			}

			got, err := svc.Create(tt.addr, tt.port)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}

			if err != nil {
				if got != nil {
					t.Fatalf("Create() got = %v, want nil", got)
				}
				return
			}

			if got.Addr != tt.addr {
				t.Fatalf("Create() got addr = %s, want %s", got.Addr, tt.addr)
			}
			if got.Port != tt.port {
				t.Fatalf("Create() got port = %d, want %d", got.Port, tt.port)
			}
		})
	}
}
