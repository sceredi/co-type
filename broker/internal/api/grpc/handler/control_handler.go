// Package handler provides the implementation for the gRPC handlers that manage incoming requests.
package handler

import (
	"context"
	"log/slog"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/broker/internal/service"
	grpc_utils "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/control"
)

// ControlHandler implements the ControlServiceServer interface and handles incoming gRPC requests for server management.
type ControlHandler struct {
	control.UnimplementedControlServiceServer
	serverSvc service.ServerService
	lobbySvc  service.LobbyService
}

// NewControlHandler creates a new instance of ControlHandler with the provided ServerService.
func NewControlHandler(serverSvc service.ServerService, lobbySvc service.LobbyService) *ControlHandler {
	return &ControlHandler{serverSvc: serverSvc, lobbySvc: lobbySvc}
}

// RegisterServer handles incoming gRPC requests to register a new server. It is idempotent: if the
// server is already known (e.g. after a broker restart) its address is updated and all supplied
// lobby IDs are re-registered.
func (h *ControlHandler) RegisterServer(_ context.Context, req *control.RegisterServerRequest) (*control.RegisterServerResponse, error) {
	h.serverSvc.Upsert(req.GetName(), req.GetHost(), int(req.GetPort()))
	for _, lobbyID := range req.GetLobbyIds() {
		h.lobbySvc.Upsert(repository.LobbyID(lobbyID), repository.ServerName(req.GetName()))
	}
	return &control.RegisterServerResponse{
		Success: true,
		Message: "Server registered successfully",
	}, nil
}

// RegisterLobby handles incoming gRPC requests to register a new lobby. It validates the request
// and upserts the lobby entry — this is idempotent and also handles the resume-game path where a
// lobby was previously hosted by a crashed server. If the lobby was pending resume, it is removed
// from the pending list.
func (h *ControlHandler) RegisterLobby(_ context.Context, req *control.RegisterLobbyRequest) (*control.RegisterLobbyResponse, error) {
	h.lobbySvc.Upsert(repository.LobbyID(req.GetLobbyId()), repository.ServerName(req.GetServerName()))
	h.lobbySvc.RemoveFromPending(repository.LobbyID(req.GetLobbyId()))
	return &control.RegisterLobbyResponse{
		Success: true,
		Message: "Lobby registered successfully",
	}, nil
}

// UnregisterLobby handles incoming gRPC requests to unregister a lobby. It validates the request and attempts to delete the lobby entry using the LobbyService.
func (h *ControlHandler) UnregisterLobby(_ context.Context, req *control.UnregisterLobbyRequest) (*control.UnregisterLobbyResponse, error) {
	err := h.lobbySvc.Delete(repository.LobbyID(req.GetLobbyId()))
	if err != nil {
		slog.Error("Failed to unregister lobby",
			slog.String("lobby_id", req.GetLobbyId()),
			slog.String("error", err.Error()),
		)
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &control.UnregisterLobbyResponse{
		Success: true,
		Message: "Lobby unregistered successfully",
	}, nil
}

// Heartbeat handles incoming heartbeat requests from game servers, updating their last-seen timestamp and load.
func (h *ControlHandler) Heartbeat(_ context.Context, req *control.HeartbeatRequest) (*control.HeartbeatResponse, error) {
	err := h.serverSvc.Heartbeat(req.GetName(), int(req.GetLoad()))
	if err != nil {
		slog.Error("Failed to process heartbeat",
			slog.String("name", req.GetName()),
			slog.String("error", err.Error()),
		)
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &control.HeartbeatResponse{Success: true}, nil
}
