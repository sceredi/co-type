// Package gateway provides an implementation of the LobbyGateway interface that interacts with the game service via gRPC.
package gateway

import (
	"context"

	"github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/lobby"
)

// LobbyGateway defines the interface for interacting with the lobby service.
type LobbyGateway interface {
	// Create creates a new lobby with the given ID and host name.
	Create(id, hostName string) (*domain.Lobby, error)
}

type lobbyGateway struct {
	ctx  context.Context
	conn lobby.LobbyServiceClient
}

// NewLobbyGateway creates a new LobbyGateway with the given gRPC connection.
func NewLobbyGateway(ctx context.Context, conn lobby.LobbyServiceClient) LobbyGateway {
	return &lobbyGateway{ctx: ctx, conn: conn}
}

func (g *lobbyGateway) Create(id, hostName string) (*domain.Lobby, error) {
	req := &lobby.CreateLobbyRequest{
		LobbyId:    id,
		PlayerName: hostName,
	}
	_, err := g.conn.CreateLobby(g.ctx, req)
	if err != nil {
		return nil, grpc.FromGRPCError(err)
	}
	return domain.NewLobby(id, domain.NewPlayer(hostName)), nil
}
