// Package gateway provides the implementation for the Gateways to connect gRPC services to the relative service.
package gateway

import (
	"github.com/sceredi/co-type/common/proto/control"
	"google.golang.org/grpc"
)

// ControlGateway defines the interface for managing control operations in the server service.
type ControlGateway interface {
	// Register registers the server with the control service using the provided name, host, and port.
	Register(name, host string, port int) error
}

type controlGateway struct {
	stream grpc.BidiStreamingClient[control.ServerEnvelope, control.BrokerEnvelope]
}

// NewControlGateway creates a new instance of ControlGateway with the provided gRPC stream and returns it.
func NewControlGateway(stream grpc.BidiStreamingClient[control.ServerEnvelope, control.BrokerEnvelope]) ControlGateway {
	return &controlGateway{stream: stream}
}

func (g *controlGateway) Register(name, host string, port int) error {
	req := &control.ServerEnvelope{
		Payload: &control.ServerEnvelope_Register{
			Register: &control.RegisterServer{
				Name: name,
				Host: host,
				Port: int64(port),
			},
		},
	}
	return g.stream.Send(req)
}
