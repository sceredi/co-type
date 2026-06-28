package gateway

import (
	"context"
	"errors"
	"testing"

	discoverypb "github.com/sceredi/co-type/common/proto/discovery"
	grpcpkg "google.golang.org/grpc"
)

// mockDiscoveryClient implements discovery.DiscoveryServiceClient for tests
type mockDiscoveryClient struct {
	resp             *discoverypb.AvailableServerResponse
	hostingLobbyResp *discoverypb.ServerHostingLobbyResponse
	err              error
}

func (m *mockDiscoveryClient) AvailableServer(ctx context.Context, in *discoverypb.AvailableServerRequest, opts ...grpcpkg.CallOption) (*discoverypb.AvailableServerResponse, error) {
	return m.resp, m.err
}

func (m *mockDiscoveryClient) ServerHostingLobby(ctx context.Context, in *discoverypb.ServerHostingLobbyRequest, opts ...grpcpkg.CallOption) (*discoverypb.ServerHostingLobbyResponse, error) {
	return m.hostingLobbyResp, m.err
}

func (m *mockDiscoveryClient) RequestResumeGame(_ context.Context, _ *discoverypb.RequestResumeGameRequest, _ ...grpcpkg.CallOption) (*discoverypb.RequestResumeGameResponse, error) {
	return nil, nil
}

// The generated client interface uses variadic grpc.CallOption; mock implements that signature.

func TestAvailableServer_Success(t *testing.T) {
	mock := &mockDiscoveryClient{resp: &discoverypb.AvailableServerResponse{Name: "s1", Addr: "127.0.0.1", Port: 9000}}
	g := &discoveryGateway{ctx: context.Background(), conn: mock}

	srv, err := g.AvailableServer()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if srv == nil {
		t.Fatal("expected server, got nil")
	}
	if srv.Name != "s1" {
		t.Fatalf("expected name s1, got %q", srv.Name)
	}
}

func TestAvailableServer_Error(t *testing.T) {
	mock := &mockDiscoveryClient{err: errors.New("gw fail")}
	g := &discoveryGateway{ctx: context.Background(), conn: mock}

	srv, err := g.AvailableServer()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if srv != nil {
		t.Fatalf("expected nil server on error, got %+v", srv)
	}
}

func TestHostByLobby_Success(t *testing.T) {
	mock := &mockDiscoveryClient{hostingLobbyResp: &discoverypb.ServerHostingLobbyResponse{Name: "s1", Addr: "127.0.0.1", Port: 9000}}
	g := &discoveryGateway{ctx: context.Background(), conn: mock}

	srv, err := g.HostByLobby("lobby1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if srv == nil {
		t.Fatal("expected server, got nil")
	}
	if srv.Name != "s1" {
		t.Fatalf("expected name s1, got %q", srv.Name)
	}
}

func TestHostByLobby_Error(t *testing.T) {
	mock := &mockDiscoveryClient{err: errors.New("gw fail")}
	g := &discoveryGateway{ctx: context.Background(), conn: mock}

	srv, err := g.HostByLobby("lobby1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if srv != nil {
		t.Fatalf("expected nil server on error, got %+v", srv)
	}
}
