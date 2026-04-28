// Package grpc implements the gRPC server for the broker service.
// It provides the necessary handlers and mappers to handle gRPC requests and responses, allowing clients to interact with the broker service using gRPC protocol.
package grpc

import (
	"errors"

	"github.com/sceredi/co-type/broker/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrServerAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	case errors.Is(err, domain.ErrServerAddrInvalid):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	case errors.Is(err, domain.ErrServerNotFound):
		return status.Errorf(codes.NotFound, "%s", err.Error())
	default:
		return status.Errorf(codes.Internal, "internal server error")
	}
}

func toGRPCMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
