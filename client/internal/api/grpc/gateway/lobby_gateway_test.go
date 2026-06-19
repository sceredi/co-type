package gateway

import (
	"context"
	"testing"

	lobbypb "github.com/sceredi/co-type/common/proto/lobby"
	grpcpkg "google.golang.org/grpc"
)

type mockLobbyClient struct {
	lastCreateReq  *lobbypb.CreateLobbyRequest
	lastJoinReq    *lobbypb.JoinLobbyRequest
	lastLeaveReq   *lobbypb.LeaveLobbyRequest
	lastEditReq    *lobbypb.EditPlayerRequest
	createResp     *lobbypb.CreateLobbyResponse
	joinResp       *lobbypb.JoinLobbyResponse
	leaveResp      *lobbypb.LeaveLobbyResponse
	editPlayerResp *lobbypb.EditPlayerResponse
	err            error
}

func (m *mockLobbyClient) CreateLobby(ctx context.Context, in *lobbypb.CreateLobbyRequest, opts ...grpcpkg.CallOption) (*lobbypb.CreateLobbyResponse, error) {
	m.lastCreateReq = in
	return m.createResp, m.err
}

func (m *mockLobbyClient) JoinLobby(ctx context.Context, in *lobbypb.JoinLobbyRequest, opts ...grpcpkg.CallOption) (*lobbypb.JoinLobbyResponse, error) {
	m.lastJoinReq = in
	return m.joinResp, m.err
}

func (m *mockLobbyClient) LeaveLobby(ctx context.Context, in *lobbypb.LeaveLobbyRequest, opts ...grpcpkg.CallOption) (*lobbypb.LeaveLobbyResponse, error) {
	m.lastLeaveReq = in
	return m.leaveResp, m.err
}

func (m *mockLobbyClient) EditPlayer(ctx context.Context, in *lobbypb.EditPlayerRequest, opts ...grpcpkg.CallOption) (*lobbypb.EditPlayerResponse, error) {
	m.lastEditReq = in
	return m.editPlayerResp, m.err
}

func (m *mockLobbyClient) Subscribe(ctx context.Context, in *lobbypb.SubscribeRequest, opts ...grpcpkg.CallOption) (grpcpkg.ServerStreamingClient[lobbypb.LobbyEvent], error) {
	return nil, nil
}

func TestCreateLobby_Success(t *testing.T) {
	mock := &mockLobbyClient{createResp: &lobbypb.CreateLobbyResponse{
		Lobby: &lobbypb.Lobby{
			Id: "ABCD",
			Host: &lobbypb.Player{
				Name: "Host",
			},
		},
	}}
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
	if mock.lastCreateReq == nil || mock.lastCreateReq.GetLobbyId() != "ABCD" {
		t.Fatalf("expected CreateLobby called with lobby id ABCD, got %+v", mock.lastCreateReq)
	}
}

func TestJoinLobby_Success(t *testing.T) {
	mock := &mockLobbyClient{joinResp: &lobbypb.JoinLobbyResponse{
		Lobby: &lobbypb.Lobby{
			Id: "ABCD",
			Players: []*lobbypb.Player{
				{Name: "Host"},
				{Name: "Player2"},
			},
			Host: &lobbypb.Player{Name: "Host"},
		},
	}}
	g := &lobbyGateway{ctx: context.Background(), conn: mock}

	l, err := g.Join("ABCD", "Player2")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l == nil {
		t.Fatal("expected lobby, got nil")
	}
	if l.ID != "ABCD" {
		t.Fatalf("expected id ABCD, got %q", l.ID)
	}
	if len(l.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(l.Players))
	}
	if mock.lastJoinReq == nil || mock.lastJoinReq.GetLobbyId() != "ABCD" {
		t.Fatalf("expected JoinLobby called with lobby id ABCD, got %+v", mock.lastJoinReq)
	}
}

func TestLeaveLobby_Success(t *testing.T) {
	mock := &mockLobbyClient{leaveResp: &lobbypb.LeaveLobbyResponse{Success: true}}
	g := &lobbyGateway{ctx: context.Background(), conn: mock}

	err := g.Leave("ABCD", "Player1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mock.lastLeaveReq == nil || mock.lastLeaveReq.GetLobbyId() != "ABCD" {
		t.Fatalf("expected LeaveLobby called with lobby id ABCD, got %+v", mock.lastLeaveReq)
	}
	if mock.lastLeaveReq.GetPlayerName() != "Player1" {
		t.Fatalf("expected player name Player1, got %q", mock.lastLeaveReq.GetPlayerName())
	}
}

func TestEditPlayer_Success(t *testing.T) {
	mock := &mockLobbyClient{editPlayerResp: &lobbypb.EditPlayerResponse{
		Lobby: &lobbypb.Lobby{
			Id: "ABCD",
			Players: []*lobbypb.Player{
				{Name: "Player1", IsReady: true, AllowedCharacters: "abc", BlockedCharacters: "xyz", CanDelete: false},
			},
			Host: &lobbypb.Player{Name: "Player1"},
		},
	}}
	g := &lobbyGateway{ctx: context.Background(), conn: mock}

	l, err := g.EditPlayer("ABCD", "Player1", true, "abc", "xyz", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l == nil {
		t.Fatal("expected lobby, got nil")
	}
	if len(l.Players) != 1 || l.Players[0].Name != "Player1" {
		t.Fatalf("unexpected players: %+v", l.Players)
	}
	if !l.Players[0].IsReady {
		t.Fatal("expected player to be ready")
	}
	if mock.lastEditReq == nil || mock.lastEditReq.GetLobbyId() != "ABCD" {
		t.Fatalf("expected EditPlayer called with lobby id ABCD, got %+v", mock.lastEditReq)
	}
	if mock.lastEditReq.GetPlayerName() != "Player1" {
		t.Fatalf("expected player name Player1, got %q", mock.lastEditReq.GetPlayerName())
	}
}
