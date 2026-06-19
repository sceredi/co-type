package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/server/internal/api/grpc/gateway"
	"github.com/sceredi/co-type/server/internal/service"
)

type mockDiscoveryGateway struct {
	name       string
	host       string
	port       int
	lobbyID    string
	serverName string
	err        error
}

func (m *mockDiscoveryGateway) RegisterServer(name, host string, port int) error {
	m.name = name
	m.host = host
	m.port = port
	return m.err
}

func (m *mockDiscoveryGateway) RegisterLobby(lobbyID, serverName string) error {
	m.lobbyID = lobbyID
	m.serverName = serverName
	return m.err
}

var _ gateway.ControlGateway = (*mockDiscoveryGateway)(nil)

func TestNewDiscoveryService(t *testing.T) {
	var svc service.DiscoveryService
	svc = service.NewDiscoveryService(&mockDiscoveryGateway{})
	if svc == nil {
		t.Fatal("NewDiscoveryService() returned nil")
	}
}

func TestDiscoveryService_Register(t *testing.T) {
	registerErr := errors.New("register failed")
	tests := []struct {
		name    string
		svcName string
		host    string
		port    int
		gateway *mockDiscoveryGateway
		wantErr error
	}{
		{
			name:    "delegates register to the gateway",
			svcName: "test-server",
			host:    "127.0.0.1",
			port:    50051,
			gateway: &mockDiscoveryGateway{},
			wantErr: nil,
		},
		{
			name:    "returns gateway error",
			svcName: "test-server",
			host:    "127.0.0.1",
			port:    50051,
			gateway: &mockDiscoveryGateway{err: registerErr},
			wantErr: registerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewDiscoveryService(tt.gateway)

			err := svc.Register(tt.svcName, tt.host, tt.port)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Register() error = %v, want %v", err, tt.wantErr)
			}

			if tt.gateway.name != tt.svcName {
				t.Fatalf("Register() gateway name = %q, want %q", tt.gateway.name, tt.svcName)
			}
			if tt.gateway.host != tt.host {
				t.Fatalf("Register() gateway host = %q, want %q", tt.gateway.host, tt.host)
			}
			if tt.gateway.port != tt.port {
				t.Fatalf("Register() gateway port = %d, want %d", tt.gateway.port, tt.port)
			}
		})
	}
}

func TestDiscoveryService_RegisterLobby(t *testing.T) {
	registerErr := errors.New("register lobby failed")
	tests := []struct {
		name       string
		lobbyID    string
		serverName string
		gateway    *mockDiscoveryGateway
		wantErr    error
	}{
		{
			name:       "delegates register lobby to the gateway",
			lobbyID:    "lobby-1",
			serverName: "test-server",
			gateway:    &mockDiscoveryGateway{},
			wantErr:    nil,
		},
		{
			name:       "returns gateway error",
			lobbyID:    "lobby-1",
			serverName: "test-server",
			gateway:    &mockDiscoveryGateway{err: registerErr},
			wantErr:    registerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewDiscoveryService(tt.gateway)

			err := svc.RegisterLobby(tt.lobbyID, tt.serverName)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("RegisterLobby() error = %v, want %v", err, tt.wantErr)
			}

			if tt.gateway.lobbyID != tt.lobbyID {
				t.Fatalf("RegisterLobby() gateway lobbyID = %q, want %q", tt.gateway.lobbyID, tt.lobbyID)
			}
			if tt.gateway.serverName != tt.serverName {
				t.Fatalf("RegisterLobby() gateway serverName = %q, want %q", tt.gateway.serverName, tt.serverName)
			}
		})
	}
}
