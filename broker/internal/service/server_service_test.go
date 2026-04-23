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
		seed    []domain.CreateServerRequest
		input   domain.CreateServerRequest
		wantErr error
	}{
		{
			name:    "creates new server",
			input:   domain.CreateServerRequest{Addr: "8.8.8.8"},
			wantErr: nil,
		},
		{
			name: "returns ErrServerAlreadyExists when server already exists",
			seed: []domain.CreateServerRequest{
				{Addr: "8.8.8.8"},
			},
			input:   domain.CreateServerRequest{Addr: "8.8.8.8"},
			wantErr: domain.ErrServerAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var svc service.ServerService
			repo := memory.NewServerRepository()
			svc = service.NewServerService(repo)
			for _, req := range tt.seed {
				if _, err := svc.Create(req); err != nil {
					t.Fatalf("seed create failed: %v", err)
				}
			}

			got, err := svc.Create(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}

			if err != nil {
				if got != nil {
					t.Fatalf("Create() got = %v, want nil", got)
				}
				return
			}

			if got.Addr != tt.input.Addr {
				t.Fatalf("Create() got = %s, want %s", got.Addr, tt.input.Addr)
			}
		})
	}
}
