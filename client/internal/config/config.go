// Package config is responsible for setting up the client.
package config

import (
	"context"

	ctgrpc "github.com/sceredi/co-type/client/internal/grpc"
	"github.com/sceredi/co-type/client/internal/service"
	"github.com/sceredi/co-type/common/proto/discovery"
	"google.golang.org/grpc"
)

// CreateDiscoveryService creates a new DiscoveryService using the provided gRPC connection.
func CreateDiscoveryService(conn *grpc.ClientConn) service.DiscoveryService {
	c := discovery.NewDiscoveryServiceClient(conn)
	gtw := ctgrpc.NewDiscoveryGateway(context.Background(), c)
	return service.NewDiscoveryService(gtw)
}
