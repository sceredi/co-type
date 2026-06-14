// Package handler provides the implementation for the gRPC handlers to manage incoming requests and interact with the relative service.
package handler

import (
	"context"

	"github.com/sceredi/co-type/common/domain"
	grpc_utils "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/lobby"
	"github.com/sceredi/co-type/server/internal/service"
	"google.golang.org/grpc/status"
)

// LobbyHandler implements the gRPC server for managing game-related operations, such as creating lobbies and handling game sessions.
type LobbyHandler struct {
	lobby.UnimplementedLobbyServiceServer
	lobbySvc service.LobbyService
}

// NewLobbyHandler creates a new instance of LobbyHandler with the provided LobbyService and returns it.
func NewLobbyHandler(lobbySvc service.LobbyService) *LobbyHandler {
	return &LobbyHandler{lobbySvc: lobbySvc}
}

// CreateLobby handles the gRPC request to create a new game lobby.
func (h *LobbyHandler) CreateLobby(_ context.Context, req *lobby.CreateLobbyRequest) (*lobby.CreateLobbyResponse, error) {
	l, err := h.lobbySvc.Create(req.LobbyId, req.PlayerName)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.CreateLobbyResponse{
		Lobby: l.Base.ToGRPCLobby(),
	}, nil
}

// JoinLobby handles the gRPC request to join an existing game lobby.
func (h *LobbyHandler) JoinLobby(_ context.Context, req *lobby.JoinLobbyRequest) (*lobby.JoinLobbyResponse, error) {
	l, err := h.lobbySvc.Join(req.LobbyId, req.PlayerName)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.JoinLobbyResponse{
		Lobby: l.Base.ToGRPCLobby(),
	}, nil
}

// Subscribe handles the gRPC subscription to lobby events.
func (h *LobbyHandler) Subscribe(req *lobby.SubscribeRequest, stream lobby.LobbyService_SubscribeServer) error {
	ctx := stream.Context()
	lobby := h.lobbySvc.Get(req.LobbyId)
	if lobby == nil {
		return grpc_utils.ToGRPCError(domain.ErrLobbyNotFound)
	}
	eventCh, ok := lobby.Subs[req.PlayerName]
	if !ok {
		return grpc_utils.ToGRPCError(domain.ErrPlayerNotInLobby)
	}
	for {
		select {
		case <-ctx.Done():
			return status.FromContextError(ctx.Err()).Err()
		case l, ok := <-eventCh:
			if !ok {
				return nil
			}
			if err := stream.Send(newLobbyEvent(l.Base)); err != nil {
				return err
			}
		}
	}
}

func newLobbyEvent(l *domain.Lobby) *lobby.LobbyEvent {
	return &lobby.LobbyEvent{
		Lobby: l.ToGRPCLobby(),
	}
}
