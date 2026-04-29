// Package config is responsible for setting up the configuration for the server application.
package config

import (
	"context"
	"log"

	"github.com/sceredi/co-type/common/proto/control"
	ctgrpc "github.com/sceredi/co-type/server/internal/grpc"
	"github.com/sceredi/co-type/server/internal/service"
	"google.golang.org/grpc"
)

// CreateDiscoveryService creates a new instance of DiscoveryService using the provided gRPC connection.
func CreateDiscoveryService(conn *grpc.ClientConn) service.DiscoveryService {
	c := control.NewControlServiceClient(conn)
	stream, err := c.Manage(context.Background())
	if err != nil {
		log.Fatalf("Error creating gRPC stream: %v", err)
	}
	gtw := ctgrpc.NewControlGateway(stream)
	return service.NewDiscoveryService(gtw)
}
