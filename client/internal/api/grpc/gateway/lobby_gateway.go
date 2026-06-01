// Package gateway provides an implementation of the LobbyGateway interface that interacts with the game service via gRPC.
package gateway

import (
	"context"
	"fmt"

	"github.com/sceredi/co-type/common/domain"
	commongrpc "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/lobby"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LobbyGateway defines the interface for interacting with the lobby service.
type LobbyGateway interface {
	// Create creates a new lobby with the given ID and host name.
	Create(id, hostName string) (*domain.Lobby, error)
	Connect(target *domain.Server) error
}

type lobbyGateway struct {
	ctx  context.Context
	conn lobby.LobbyServiceClient
}

// NewLobbyGateway creates a new LobbyGateway with the given gRPC connection.
func NewLobbyGateway(ctx context.Context) LobbyGateway {
	return &lobbyGateway{ctx: ctx}
}

func (g *lobbyGateway) Create(id, hostName string) (*domain.Lobby, error) {
	req := &lobby.CreateLobbyRequest{
		LobbyId:    id,
		PlayerName: hostName,
	}
	_, err := g.conn.CreateLobby(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	return domain.NewLobby(id, domain.NewPlayer(hostName)), nil
}

func (g *lobbyGateway) Connect(target *domain.Server) error {
	addr := fmt.Sprintf("%s:%d", target.Addr, target.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return commongrpc.FromGRPCError(err)
	}
	g.conn = lobby.NewLobbyServiceClient(conn)
	return nil
}
