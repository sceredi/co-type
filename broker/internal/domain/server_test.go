package domain_test

import (
	"errors"
	"testing"

	"github.com/sceredi/co-type/broker/internal/domain"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr error
	}{
		{
			name:    "creates a new server",
			addr:    "8.8.8.8",
			wantErr: nil,
		},
		{
			name:    "returns ErrServerAddrInvalid when addr is empty",
			addr:    "",
			wantErr: domain.ErrServerAddrInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewServer(tt.addr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewServer() failed: %v", err)
			}
			if err != nil {
				if got != nil {
					t.Fatalf("NewServer() got = %v, want nil", got)
				}
				return
			}

			if got.Addr != tt.addr {
				t.Fatalf("NewServer() got addr = %s, want %s", got.Addr, tt.addr)
			}
		})
	}
}
