package domain_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/common/domain"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name     string
		hostName string
		addr     string
		port     int32
		wantErr  error
	}{
		{
			name:     "creates a new server",
			hostName: "testserver1",
			addr:     "8.8.8.8",
			port:     8080,
			wantErr:  nil,
		},
		{
			name:     "returns ErrServerNameInvalid when hostName is empty",
			hostName: "",
			addr:     "8.8.8.8",
			port:     8080,
			wantErr:  domain.ErrServerNameInvalid,
		},
		{
			name:     "returns ErrServerAddrInvalid when addr is empty",
			hostName: "testserver1",
			addr:     "",
			port:     8080,
			wantErr:  domain.ErrServerAddrInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewServer(tt.hostName, tt.addr, tt.port)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewServer() failed: %v", err)
			}
			if err != nil {
				if got != nil {
					t.Fatalf("NewServer() got = %v, want nil", got)
				}
				return
			}

			if got.Name != tt.hostName {
				t.Fatalf("NewServer() got name = %s, want %s", got.Name, tt.hostName)
			}
			if got.Addr != tt.addr {
				t.Fatalf("NewServer() got addr = %s, want %s", got.Addr, tt.addr)
			}
			if got.Port != tt.port {
				t.Fatalf("NewServer() got port = %d, want %d", got.Port, tt.port)
			}
			if got.Load != 0 {
				t.Fatalf("NewServer() got load = %d, want %d", got.Load, 0)
			}
		})
	}
}
