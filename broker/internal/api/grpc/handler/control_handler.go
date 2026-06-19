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

// RegisterServer handles incoming gRPC requests to register a new server. It validates the request and attempts to create a new server entry using the ServerService.
func (h *ControlHandler) RegisterServer(_ context.Context, req *control.RegisterServerRequest) (*control.RegisterServerResponse, error) {
	_, err := h.serverSvc.Create(req.GetName(), req.GetHost(), int(req.GetPort()))
	if err != nil {
		slog.Error("Failed to register server",
			slog.String("name", req.GetName()),
			slog.String("host", req.GetHost()),
			slog.Int("port", int(req.GetPort())),
			slog.String("error", err.Error()),
		)
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &control.RegisterServerResponse{
		Success: true,
		Message: "Server registered successfully",
	}, nil
}

// RegisterLobby handles incoming gRPC requests to register a new lobby. It validates the request and attempts to create a new lobby entry using the LobbyService.
func (h *ControlHandler) RegisterLobby(_ context.Context, req *control.RegisterLobbyRequest) (*control.RegisterLobbyResponse, error) {
	err := h.lobbySvc.Create(repository.LobbyID(req.GetLobbyId()), repository.ServerName(req.GetServerName()))
	if err != nil {
		slog.Error("Failed to register lobby",
			slog.String("lobby_id", req.GetLobbyId()),
			slog.String("server_name", req.GetServerName()),
			slog.String("error", err.Error()),
		)
		return nil, grpc_utils.ToGRPCError(err)
	}
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
