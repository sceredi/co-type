// Package handler provides the implementation for the gRPC handlers to manage incoming requests and interact with the relative service.
package handler

import (
	"context"

	grpc_utils "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/game"
	"github.com/sceredi/co-type/server/internal/service"
)

// GameHandler implements the gRPC server for managing game-related operations, such as creating lobbies and handling game sessions.
type GameHandler struct {
	game.UnimplementedGameServiceServer
	lobbySvc service.LobbyService
}

// NewGameHandler creates a new instance of GameHandler with the provided LobbyService and returns it.
func NewGameHandler(lobbySvc service.LobbyService) *GameHandler {
	return &GameHandler{lobbySvc: lobbySvc}
}

// CreateLobby handles the gRPC request to create a new game lobby.
func (h *GameHandler) CreateLobby(_ context.Context, req *game.CreateLobbyRequest) (*game.CreateLobbyResponse, error) {
	_, err := h.lobbySvc.Create(req.LobbyId, req.PlayerName)
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &game.CreateLobbyResponse{}, nil
}
