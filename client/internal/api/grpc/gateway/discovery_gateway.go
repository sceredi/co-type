// Package gateway contains all the gateway that the client uses to connect to the broker and game servers.
package gateway

import (
	"context"
	"log/slog"

	"github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/discovery"
)

// DiscoveryGateway manages the interaction with the broker service.
type DiscoveryGateway interface {
	// AvailableServer returns either a server that is available to manage a lobby or an error if something goes wrong.
	AvailableServer() (*domain.Server, error)
	// HostByLobby returns the server that is hosting the lobby with the given code or an error if something goes wrong.
	HostByLobby(lobbyCode string) (*domain.Server, error)
	// RequestResumeGame asks the broker to assign a server for resuming the given crashed lobby.
	// All players in the crashed game should call this; the broker returns the same server for each.
	RequestResumeGame(lobbyID string) (*domain.Server, error)
}

type discoveryGateway struct {
	ctx  context.Context
	conn discovery.DiscoveryServiceClient
}

// NewDiscoveryGateway creates a new DiscoveryGateway.
func NewDiscoveryGateway(ctx context.Context, conn discovery.DiscoveryServiceClient) DiscoveryGateway {
	return &discoveryGateway{ctx: ctx, conn: conn}
}

func (g *discoveryGateway) AvailableServer() (*domain.Server, error) {
	req := &discovery.AvailableServerRequest{}
	srv, err := g.conn.AvailableServer(g.ctx, req)
	slog.Debug("Requested available server", "request", req, "response", srv)
	if err != nil {
		return nil, grpc.FromGRPCError(err)
	}
	return domain.NewServer(srv.Name, srv.Addr, int(srv.Port))
}

func (g *discoveryGateway) HostByLobby(lobbyCode string) (*domain.Server, error) {
	req := &discovery.ServerHostingLobbyRequest{
		LobbyId: lobbyCode,
	}
	srv, err := g.conn.ServerHostingLobby(g.ctx, req)
	slog.Debug("Requested server hosting lobby", "request", req, "response", srv)
	if err != nil {
		return nil, grpc.FromGRPCError(err)
	}
	return domain.NewServer(srv.Name, srv.Addr, int(srv.Port))
}

func (g *discoveryGateway) RequestResumeGame(lobbyID string) (*domain.Server, error) {
	req := &discovery.RequestResumeGameRequest{LobbyId: lobbyID}
	srv, err := g.conn.RequestResumeGame(g.ctx, req)
	slog.Debug("Requested resume game server", "request", req, "response", srv)
	if err != nil {
		return nil, grpc.FromGRPCError(err)
	}
	return domain.NewServer(srv.Name, srv.Addr, int(srv.Port))
}
