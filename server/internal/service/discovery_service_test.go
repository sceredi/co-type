package service_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/server/internal/grpc"
	"github.com/sceredi/co-type/server/internal/service"
)

type mockDiscoveryGateway struct {
	name string
	host string
	port int32
	err  error
}

func (m *mockDiscoveryGateway) Register(name, host string, port int32) error {
	m.name = name
	m.host = host
	m.port = port
	return m.err
}

var _ grpc.ControlGateway = (*mockDiscoveryGateway)(nil)

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
		port    int32
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
