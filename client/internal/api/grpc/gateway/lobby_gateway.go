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
	Join(id, playerName string) (*domain.Lobby, error)
	Leave(id, playerName string) error
	EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error)
	Ready(lobbyID, playerName string) (*domain.Lobby, error)
	Connect(target *domain.Server) error
	// Subscribe subscribes to lobby events and returns a channel that receives updated lobby state.
	Subscribe(lobbyID, playerName string) (<-chan *domain.Lobby, error)
	// SendKeyPress sends a validated key press to the server.
	SendKeyPress(lobbyID, playerName, key string, isBackspace bool) (*domain.Lobby, error)
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
	res, err := g.conn.CreateLobby(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	return domain.NewLobbyFromGRPC(res.GetLobby()), nil
}

func (g *lobbyGateway) Join(id, playerName string) (*domain.Lobby, error) {
	req := &lobby.JoinLobbyRequest{
		LobbyId:    id,
		PlayerName: playerName,
	}
	res, err := g.conn.JoinLobby(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	return domain.NewLobbyFromGRPC(res.GetLobby()), nil
}

func (g *lobbyGateway) Leave(id, playerName string) error {
	req := &lobby.LeaveLobbyRequest{
		LobbyId:    id,
		PlayerName: playerName,
	}
	_, err := g.conn.LeaveLobby(g.ctx, req)
	if err != nil {
		return commongrpc.FromGRPCError(err)
	}
	return nil
}

func (g *lobbyGateway) EditPlayer(lobbyID, playerName string, isReady bool, allowedCharacters, blockedCharacters string, canDelete bool) (*domain.Lobby, error) {
	req := &lobby.EditPlayerRequest{
		LobbyId:           lobbyID,
		PlayerName:        playerName,
		IsReady:           isReady,
		AllowedCharacters: allowedCharacters,
		BlockedCharacters: blockedCharacters,
		CanDelete:         canDelete,
	}
	res, err := g.conn.EditPlayer(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	return domain.NewLobbyFromGRPC(res.GetLobby()), nil
}

func (g *lobbyGateway) Ready(lobbyID, playerName string) (*domain.Lobby, error) {
	req := &lobby.ReadyPlayerRequest{
		LobbyId:    lobbyID,
		PlayerName: playerName,
	}
	res, err := g.conn.ReadyPlayer(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	return domain.NewLobbyFromGRPC(res.GetLobby()), nil
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

func (g *lobbyGateway) Subscribe(lobbyID, playerName string) (<-chan *domain.Lobby, error) {
	req := &lobby.SubscribeRequest{
		LobbyId:    lobbyID,
		PlayerName: playerName,
	}
	stream, err := g.conn.Subscribe(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	ch := make(chan *domain.Lobby, 64)
	go func() {
		defer close(ch)
		for {
			event, err := stream.Recv()
			if err != nil {
				return
			}
			ch <- domain.NewLobbyFromGRPC(event.GetLobby())
		}
	}()
	return ch, nil
}

func (g *lobbyGateway) SendKeyPress(lobbyID, playerName, key string, isBackspace bool) (*domain.Lobby, error) {
	req := &lobby.SendKeyPressRequest{
		LobbyId:     lobbyID,
		PlayerName:  playerName,
		Key:         key,
		IsBackspace: isBackspace,
	}
	res, err := g.conn.SendKeyPress(g.ctx, req)
	if err != nil {
		return nil, commongrpc.FromGRPCError(err)
	}
	return domain.NewLobbyFromGRPC(res.GetLobby()), nil
}
