// Package gateway provides the implementation for the Gateways to connect gRPC services to the relative service.
package gateway

import (
	"context"
	"errors"

	commongrpc "github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/control"
)

// ControlGateway defines the interface for managing control operations in the server service.
type ControlGateway interface {
	// RegisterServer registers the server with the control service using the provided name, host, port, and
	// the list of lobby IDs currently managed by this server (empty on first startup).
	RegisterServer(name, host string, port int, lobbyIDs []string) error

	RegisterLobby(lobbyID, serverName string) error

	UnregisterLobby(lobbyID string) error

	// Heartbeat sends a heartbeat to the broker with the server's current load.
	Heartbeat(name string, load int) error
}

type controlGateway struct {
	ctx  context.Context
	conn control.ControlServiceClient
}

// NewControlGateway creates a new instance of ControlGateway with the provided gRPC stream and returns it.
func NewControlGateway(ctx context.Context, conn control.ControlServiceClient) ControlGateway {
	return &controlGateway{ctx: ctx, conn: conn}
}

func (g *controlGateway) RegisterServer(name, host string, port int, lobbyIDs []string) error {
	req := &control.RegisterServerRequest{
		Name:     name,
		Host:     host,
		Port:     int64(port),
		LobbyIds: lobbyIDs,
	}
	res, err := g.conn.RegisterServer(g.ctx, req)
	if err != nil {
		return commongrpc.FromGRPCError(err)
	}
	if !res.Success {
		return errors.New("failed to register server: " + res.Message)
	}
	return nil
}

func (g *controlGateway) RegisterLobby(lobbyID, serverName string) error {
	req := &control.RegisterLobbyRequest{
		LobbyId:    lobbyID,
		ServerName: serverName,
	}
	res, err := g.conn.RegisterLobby(g.ctx, req)
	if err != nil {
		return commongrpc.FromGRPCError(err)
	}
	if !res.Success {
		return errors.New("failed to register lobby: " + res.Message)
	}
	return nil
}

func (g *controlGateway) UnregisterLobby(lobbyID string) error {
	req := &control.UnregisterLobbyRequest{
		LobbyId: lobbyID,
	}
	res, err := g.conn.UnregisterLobby(g.ctx, req)
	if err != nil {
		return commongrpc.FromGRPCError(err)
	}
	if !res.Success {
		return errors.New("failed to unregister lobby: " + res.Message)
	}
	return nil
}

func (g *controlGateway) Heartbeat(name string, load int) error {
	req := &control.HeartbeatRequest{
		Name: name,
		Load: int64(load),
	}
	res, err := g.conn.Heartbeat(g.ctx, req)
	if err != nil {
		return commongrpc.FromGRPCError(err)
	}
	if !res.Success {
		return errors.New("heartbeat rejected by broker")
	}
	return nil
}
