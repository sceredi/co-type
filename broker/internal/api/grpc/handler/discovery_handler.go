package handler

import (
	"context"

	"github.com/sceredi/co-type/broker/internal/repository"
	"github.com/sceredi/co-type/broker/internal/service"
	grpc_utils "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/discovery"
)

// DiscoveryHandler implements the gRPC server for the DicoveryService.
type DiscoveryHandler struct {
	discovery.UnimplementedDiscoveryServiceServer
	serverSvc service.ServerService
	lobbySvc  service.LobbyService
}

// NewDiscoveryHandler creates a new instance of DiscoveryHandler with the provided ServerService.
func NewDiscoveryHandler(serverSvc service.ServerService, lobbySvc service.LobbyService) *DiscoveryHandler {
	return &DiscoveryHandler{serverSvc: serverSvc, lobbySvc: lobbySvc}
}

// AvailableServer handles incoming gRPC requests to find the server with the lowest load.
// It returns the server information or an error if there are no available servers.
func (h *DiscoveryHandler) AvailableServer(_ context.Context, _ *discovery.AvailableServerRequest) (*discovery.AvailableServerResponse, error) {
	server, err := h.serverSvc.LowestLoad()
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &discovery.AvailableServerResponse{
		Name: server.Name,
		Addr: server.Addr,
		Port: int64(server.Port),
	}, nil
}

// ServerHostingLobby handles incoming gRPC requests to find the server hosting a specific lobby.
func (h *DiscoveryHandler) ServerHostingLobby(_ context.Context, req *discovery.ServerHostingLobbyRequest) (*discovery.ServerHostingLobbyResponse, error) {
	server, err := h.lobbySvc.Get(repository.LobbyID(req.LobbyId))
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &discovery.ServerHostingLobbyResponse{
		Name: server.Name,
		Addr: server.Addr,
		Port: int64(server.Port),
	}, nil
}

// RequestResumeGame handles incoming gRPC requests to resume a crashed game. It assigns (or
// reuses) a server for the lobby and returns its address.
func (h *DiscoveryHandler) RequestResumeGame(_ context.Context, req *discovery.RequestResumeGameRequest) (*discovery.RequestResumeGameResponse, error) {
	server, err := h.lobbySvc.RequestResumeGame(repository.LobbyID(req.LobbyId))
	if err != nil {
		return nil, grpc_utils.ToGRPCError(err)
	}
	return &discovery.RequestResumeGameResponse{
		Name: server.Name,
		Addr: server.Addr,
		Port: int64(server.Port),
	}, nil
}
