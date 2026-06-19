package service

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
)

// mockDiscoveryGateway is a mock implementation of DiscoveryGateway for testing
type mockDiscoveryGateway struct {
	server *domain.Server
	err    error
}

func (m *mockDiscoveryGateway) AvailableServer() (*domain.Server, error) {
	return m.server, m.err
}

func (m *mockDiscoveryGateway) HostByLobby(lobbyCode string) (*domain.Server, error) {
	return m.server, m.err
}

func TestNewDiscoveryService(t *testing.T) {
	gtw := &mockDiscoveryGateway{}
	svc := NewDiscoveryService(gtw)

	if svc == nil {
		t.Fatal("expected NewDiscoveryService to return non-nil service")
	}
}

func TestGetAvailableServer(t *testing.T) {
	tests := []struct {
		name     string
		gateway  *mockDiscoveryGateway
		wantErr  bool
		wantNil  bool
		wantName string
	}{
		{
			name: "delegates_to_gateway_on_success",
			gateway: &mockDiscoveryGateway{
				server: &domain.Server{Name: "server1", Addr: "10.0.0.1", Port: 8080},
			},
			wantErr:  false,
			wantNil:  false,
			wantName: "server1",
		},
		{
			name: "propagates_gateway_error",
			gateway: &mockDiscoveryGateway{
				err: errors.New("gateway unavailable"),
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "returns_nil_server_on_error",
			gateway: &mockDiscoveryGateway{
				err: errors.New("connection failed"),
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDiscoveryService(tt.gateway)
			srv, err := svc.GetAvailableServer()

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got error=%v", tt.wantErr, err != nil)
			}

			if (srv == nil) != tt.wantNil {
				t.Fatalf("expected nil=%v, got nil=%v", tt.wantNil, srv == nil)
			}

			if !tt.wantNil && srv.Name != tt.wantName {
				t.Fatalf("expected name=%q, got name=%q", tt.wantName, srv.Name)
			}
		})
	}
}

func TestGetHostByLobby(t *testing.T) {
	tests := []struct {
		name     string
		gateway  *mockDiscoveryGateway
		wantErr  bool
		wantNil  bool
		wantName string
	}{
		{
			name: "delegates_to_gateway_on_success",
			gateway: &mockDiscoveryGateway{
				server: &domain.Server{Name: "server1", Addr: "10.0.0.1", Port: 8080},
			},
			wantErr:  false,
			wantNil:  false,
			wantName: "server1",
		},
		{
			name: "propagates_gateway_error",
			gateway: &mockDiscoveryGateway{
				err: errors.New("gateway unavailable"),
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDiscoveryService(tt.gateway)
			srv, err := svc.GetHostByLobby("lobby1")

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got error=%v", tt.wantErr, err != nil)
			}

			if (srv == nil) != tt.wantNil {
				t.Fatalf("expected nil=%v, got nil=%v", tt.wantNil, srv == nil)
			}

			if !tt.wantNil && srv.Name != tt.wantName {
				t.Fatalf("expected name=%q, got name=%q", tt.wantName, srv.Name)
			}
		})
	}
}
