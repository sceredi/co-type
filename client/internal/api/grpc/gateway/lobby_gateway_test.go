package gateway

import (
	"context"
	"testing"

	lobbypb "github.com/sceredi/co-type/common/proto/lobby"
	grpcpkg "google.golang.org/grpc"
)

type mockLobbyClient struct {
	lastReq *lobbypb.CreateLobbyRequest
	resp    *lobbypb.CreateLobbyResponse
	err     error
}

func (m *mockLobbyClient) CreateLobby(ctx context.Context, in *lobbypb.CreateLobbyRequest, opts ...grpcpkg.CallOption) (*lobbypb.CreateLobbyResponse, error) {
	m.lastReq = in
	return m.resp, m.err
}

func (m *mockLobbyClient) Subscribe(ctx context.Context, in *lobbypb.SubscribeRequest, opts ...grpcpkg.CallOption) (grpcpkg.ServerStreamingClient[lobbypb.LobbyEvent], error) {
	return nil, nil
}

func TestCreateLobby_Success(t *testing.T) {
	mock := &mockLobbyClient{resp: &lobbypb.CreateLobbyResponse{}}
	g := &lobbyGateway{ctx: context.Background(), conn: mock}

	l, err := g.Create("ABCD", "Host")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l == nil {
		t.Fatal("expected lobby, got nil")
	}
	if l.ID != "ABCD" {
		t.Fatalf("expected id ABCD, got %q", l.ID)
	}
	if l.Host.Name != "Host" {
		t.Fatalf("expected host Host, got %q", l.Host.Name)
	}
	if mock.lastReq == nil || mock.lastReq.GetLobbyId() != "ABCD" {
		t.Fatalf("expected CreateLobby called with lobby id ABCD, got %+v", mock.lastReq)
	}
}
