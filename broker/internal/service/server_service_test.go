package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/broker/internal/repository/memory"
	"github.com/sceredi/co-type/broker/internal/service"
	"github.com/sceredi/co-type/common/domain"
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
		name     string
		seed     []*domain.Server
		hostName string
		addr     string
		port     int
		wantErr  error
	}{
		{
			name:     "creates new server",
			hostName: "test_server",
			addr:     "8.8.8.8",
			port:     8080,
			wantErr:  nil,
		},
		{
			name: "returns ErrServerAlreadyExists when server already exists",
			seed: []*domain.Server{
				{Name: "test_server", Addr: "8.8.8.8", Port: 8080},
			},
			hostName: "test_server",
			addr:     "10.0.0.1",
			port:     9090,
			wantErr:  domain.ErrServerAlreadyExists,
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

			got, err := svc.Create(tt.hostName, tt.addr, tt.port)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}

			if err != nil {
				if got != nil {
					t.Fatalf("Create() got = %v, want nil", got)
				}
				return
			}

			if got.Name != tt.hostName {
				t.Fatalf("NewServer() got name = %s, want %s", got.Name, tt.hostName)
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

func TestServerService_LowestLoad(t *testing.T) {
	tests := []struct {
		name      string
		seed      []*domain.Server
		wantName  string
		wantLoad  int
		wantError error
	}{
		{
			name:      "returns ErrNoAvailableServers when repository is empty",
			seed:      nil,
			wantError: domain.ErrNoAvailableServers,
		},
		{
			name: "returns server with lowest load",
			seed: []*domain.Server{
				{Name: "server-high", Addr: "10.0.0.1", Port: 8080, Load: 10},
				{Name: "server-low", Addr: "10.0.0.2", Port: 8081, Load: 2},
				{Name: "server-mid", Addr: "10.0.0.3", Port: 8082, Load: 5},
			},
			wantName:  "server-low",
			wantLoad:  2,
			wantError: nil,
		},
		{
			name: "returns first server when loads are equal",
			seed: []*domain.Server{
				{Name: "server-first", Addr: "10.0.0.4", Port: 8080, Load: 3},
				{Name: "server-second", Addr: "10.0.0.5", Port: 8081, Load: 3},
			},
			wantName:  "server-first",
			wantLoad:  3,
			wantError: nil,
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

			got, err := svc.LowestLoad()
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("LowestLoad() error = %v, want %v", err, tt.wantError)
			}

			if err != nil {
				if got != nil {
					t.Fatalf("LowestLoad() got = %v, want nil", got)
				}
				return
			}

			if got.Name != tt.wantName {
				t.Fatalf("LowestLoad() got name = %s, want %s", got.Name, tt.wantName)
			}
			if got.Load != tt.wantLoad {
				t.Fatalf("LowestLoad() got load = %d, want %d", got.Load, tt.wantLoad)
			}
		})
	}
}
