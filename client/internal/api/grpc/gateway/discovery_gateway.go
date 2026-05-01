// Package gateway contains all the gateway that the client uses to connect to the broker and game servers.
package gateway

import (
	"context"

	"github.com/sceredi/co-type/common/domain"
	"github.com/sceredi/co-type/common/grpc"
	"github.com/sceredi/co-type/common/proto/discovery"
)

// DiscoveryGateway manages the interaction with the broker service.
type DiscoveryGateway interface {
	// AvailableServer returns either a server that is available to manage a lobby or an error if something goes wrong.
	AvailableServer() (*domain.Server, error)
}

type discoveryGateway struct {
	ctx  context.Context
	conn discovery.DiscoveryServiceClient
}

// NewDiscoveryGateway creates a new DiscoveryGateway.
func NewDiscoveryGateway(ctx context.Context, conn discovery.DiscoveryServiceClient) DiscoveryGateway {
	return &discoveryGateway{ctx: ctx, conn: conn}
}

func (g *discoveryGateway) AvailableServer() (*domain.Server, error) {
	req := &discovery.AvailableServerRequest{}
	srv, err := g.conn.AvailableServer(g.ctx, req)
	if err != nil {
		return nil, grpc.FromGRPCError(err)
	}
	return domain.NewServer(srv.Name, srv.Addr, srv.Port)
}
