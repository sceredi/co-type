// Package handler provides the implementation for the gRPC handlers to manage incoming requests and interact with the relative service.
package handler

import (
	"errors"
	"log/slog"

	grpc_utils "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/game"
	"github.com/sceredi/co-type/server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Manage handles the bidirectional streaming for managing game-related operations.
// It listens for incoming messages from the client, processes them based on their type, and sends appropriate responses back to the client.
// The method continues to listen for messages until the stream is closed by the client or an error occurs.
func (h *GameHandler) Manage(stream game.GameService_ManageServer) error {
	for {
		env, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled {
				slog.DebugContext(stream.Context(), "Stream canceled by client")
			} else {
				slog.ErrorContext(stream.Context(), "Failed to receive message from stream",
					slog.String("error", err.Error()),
				)
				return err
			}
		}
		slog.DebugContext(stream.Context(), "Received message from stream",
			slog.Any("payload", env.GetPayload()),
		)
		switch msg := env.GetPayload().(type) {
		case *game.ServerEnvelope_Create:
			err = h.manageCreateReq(msg.Create, stream)
			if err != nil {
				return err
			}
		default:
			slog.ErrorContext(stream.Context(), "Received unknown message type",
				slog.Any("payload", env.GetPayload()),
			)
			return grpc_utils.ToGRPCError(errors.New("unknown message type"))
		}
	}
}

func (h *GameHandler) manageCreateReq(msg *game.CreateLobby, stream game.GameService_ManageServer) error {
	_, err := h.lobbySvc.Create(msg.GetLobbyId(), msg.GetPlayerName())
	if err != nil {
		slog.ErrorContext(stream.Context(), "Failed to create lobby",
			slog.String("id", msg.GetLobbyId()),
			slog.String("playerName", msg.GetPlayerName()),
		)
	}
	err = stream.Send(
		&game.ClientEnvelope{
			Payload: &game.ClientEnvelope_CreateAck{
				CreateAck: &game.CreateLobbyAck{
					Success: err == nil,
					Message: grpc_utils.ToGRPCMessage(err),
				},
			},
		},
	)
	if err != nil {
		slog.ErrorContext(stream.Context(), "Stream unexpectedly terminated")
	}
	return err
}
