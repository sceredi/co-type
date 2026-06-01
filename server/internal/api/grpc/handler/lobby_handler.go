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
	_, err := h.lobbySvc.Create(req.LobbyId, req.PlayerName)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.CreateLobbyResponse{}, nil
}

// Subscribe handles the gRPC subscription to lobby events.
func (h *LobbyHandler) Subscribe(req *lobby.SubscribeRequest, stream lobby.LobbyService_SubscribeServer) error {
	ctx := stream.Context()
	lobby := h.lobbySvc.Get(req.LobbyId)
	if lobby == nil {
		return domain.ErrLobbyNotFound
	}
	eventCh, ok := lobby.Subs[req.PlayerName]
	if !ok {
		return domain.ErrPlayerNotInLobby
	}
	for {
		select {
		case <-ctx.Done():
			return status.FromContextError(ctx.Err()).Err()
		case ev, ok := <-eventCh:
			if !ok {
				return nil
			}
			if err := stream.Send(toProtoEvent(ev)); err != nil {
				return err
			}
		}
	}
}

func toProtoEvent(e domain.LobbyEvent) *lobby.LobbyEvent {
	switch e := e.(type) {
	case domain.PlayerJoin:
		return &lobby.LobbyEvent{
			Event: &lobby.LobbyEvent_PlayerJoin{
				PlayerJoin: &lobby.PlayerJoin{
					PlayerName: e.PlayerName,
				},
			},
		}
	case domain.PlayerLeave:
		return &lobby.LobbyEvent{
			Event: &lobby.LobbyEvent_PlayerLeave{
				PlayerLeave: &lobby.PlayerLeave{
					PlayerName: e.PlayerName,
				},
			},
		}
	case domain.PlayerReady:
		return &lobby.LobbyEvent{
			Event: &lobby.LobbyEvent_PlayerReady{
				PlayerReady: &lobby.PlayerReady{
					PlayerName: e.PlayerName,
					IsReady:    e.IsReady,
				},
			},
		}
	case domain.PlayerEdit:
		return &lobby.LobbyEvent{
			Event: &lobby.LobbyEvent_PlayerEdit{
				PlayerEdit: &lobby.PlayerEdit{
					PlayerName:        e.PlayerName,
					AllowedCharacters: e.AllowedCharacters,
					BlockedCharacters: e.BlockedCharacters,
					CanDelete:         e.CanDelete,
				},
			},
		}
	default:
		return nil
	}
}
