// Package handler provides the implementation for the gRPC handlers to manage incoming requests and interact with the relative service.
package handler

import (
	"context"
	"log/slog"

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

// LeaveLobby handles the gRPC request to leave a game lobby.
func (h *LobbyHandler) LeaveLobby(_ context.Context, req *lobby.LeaveLobbyRequest) (*lobby.LeaveLobbyResponse, error) {
	err := h.lobbySvc.Leave(req.LobbyId, req.PlayerName)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.LeaveLobbyResponse{
		Success: true,
		Message: "Left lobby successfully",
	}, nil
}

// EditPlayer handles the gRPC request to edit a player's settings in a lobby.
func (h *LobbyHandler) EditPlayer(_ context.Context, req *lobby.EditPlayerRequest) (*lobby.EditPlayerResponse, error) {
	l, err := h.lobbySvc.EditPlayer(req.LobbyId, req.PlayerName, req.IsReady, req.AllowedCharacters, req.BlockedCharacters, req.CanDelete)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.EditPlayerResponse{
		Lobby: l.Base.ToGRPCLobby(),
	}, nil
}

// ReadyPlayer handles the gRPC request to toggle a player's ready status in a lobby.
func (h *LobbyHandler) ReadyPlayer(_ context.Context, req *lobby.ReadyPlayerRequest) (*lobby.ReadyPlayerResponse, error) {
	l, err := h.lobbySvc.Ready(req.LobbyId, req.PlayerName)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.ReadyPlayerResponse{
		Lobby: l.Base.ToGRPCLobby(),
	}, nil
}

// Subscribe handles the gRPC subscription to lobby events.
func (h *LobbyHandler) Subscribe(req *lobby.SubscribeRequest, stream lobby.LobbyService_SubscribeServer) error {
	ctx := stream.Context()
	l := h.lobbySvc.Get(req.LobbyId)
	if l == nil {
		return grpc_utils.ToGRPCError(domain.ErrLobbyNotFound)
	}
	eventCh, ok := l.Subs[req.PlayerName]
	if !ok {
		return grpc_utils.ToGRPCError(domain.ErrPlayerNotInLobby)
	}
	for {
		select {
		case <-ctx.Done():
			h.handleSubscriberDisconnect(ctx, req.LobbyId, req.PlayerName)
			return status.FromContextError(ctx.Err()).Err()
		case event, ok := <-eventCh:
			if !ok {
				return nil
			}
			if err := stream.Send(newLobbyEvent(event.Base)); err != nil {
				h.handleSubscriberDisconnect(ctx, req.LobbyId, req.PlayerName)
				return err
			}
		}
	}
}

// handleSubscriberDisconnect cleans up after a subscriber's connection is lost.
// In lobby phase the player is removed; in game phase the game is paused and the
// player is kept so they can reconnect.
func (h *LobbyHandler) handleSubscriberDisconnect(ctx context.Context, lobbyID, playerName string) {
	current := h.lobbySvc.Get(lobbyID)
	if current == nil {
		return
	}
	switch current.Base.Status {
	case domain.LobbyWaitingForPlayers:
		if err := h.lobbySvc.Leave(lobbyID, playerName); err != nil {
			slog.ErrorContext(ctx, "Failed to remove disconnected player from lobby",
				slog.String("lobbyID", lobbyID),
				slog.String("player", playerName),
				slog.String("error", err.Error()),
			)
		}
	case domain.LobbyPlaying, domain.LobbyPaused:
		if err := h.lobbySvc.PlayerDisconnected(lobbyID, playerName); err != nil {
			slog.ErrorContext(ctx, "Failed to mark player as disconnected",
				slog.String("lobbyID", lobbyID),
				slog.String("player", playerName),
				slog.String("error", err.Error()),
			)
		}
	default:
		slog.DebugContext(ctx, "Subscriber disconnected in unhandled lobby state, ignoring",
			slog.String("lobbyID", lobbyID),
			slog.String("player", playerName),
			slog.Int("status", int(current.Base.Status)),
		)
	}
}

func newLobbyEvent(l *domain.Lobby) *lobby.LobbyEvent {
	return &lobby.LobbyEvent{
		Lobby: l.ToGRPCLobby(),
	}
}

// SendKeyPress handles the gRPC request to send a key press to the game.
func (h *LobbyHandler) SendKeyPress(_ context.Context, req *lobby.SendKeyPressRequest) (*lobby.SendKeyPressResponse, error) {
	l, err := h.lobbySvc.SendKeyPress(req.LobbyId, req.PlayerName, req.Key, req.IsBackspace)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &lobby.SendKeyPressResponse{
		Lobby: l.Base.ToGRPCLobby(),
	}, nil
}
