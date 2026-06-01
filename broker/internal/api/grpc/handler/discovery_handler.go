package handler

import (
	"context"

	"github.com/sceredi/co-type/broker/internal/service"
	grpc_utils "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/discovery"
)

// DiscoveryHandler implements the gRPC server for the DicoveryService.
type DiscoveryHandler struct {
	discovery.UnimplementedDiscoveryServiceServer
	serverSvc service.ServerService
}

// NewDiscoveryHandler creates a new instance of DiscoveryHandler with the provided ServerService.
func NewDiscoveryHandler(serverSvc service.ServerService) *DiscoveryHandler {
	return &DiscoveryHandler{serverSvc: serverSvc}
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
