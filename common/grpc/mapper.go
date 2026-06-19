// Package grpc provides helpers for gRPC communication.
package grpc

import (
	"errors"
	"log"

	"github.com/sceredi/co-type/common/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError maps domain errors to gRPC status errors, allowing for consistent error handling in gRPC responses.
func ToGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrServerAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	case errors.Is(err, domain.ErrServerAddrInvalid):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	case errors.Is(err, domain.ErrServerNotFound):
		return status.Errorf(codes.NotFound, "%s", err.Error())
	case errors.Is(err, domain.ErrNoAvailableServers):
		return status.Errorf(codes.Unavailable, "%s", err.Error())
	case errors.Is(err, domain.ErrLobbyAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	case errors.Is(err, domain.ErrLobbyNotFound):
		return status.Errorf(codes.NotFound, "%s", err.Error())
	case errors.Is(err, domain.ErrPlayerNotInLobby):
		return status.Errorf(codes.NotFound, "%s", err.Error())
	case errors.Is(err, domain.ErrPlayerAlreadyInLobby):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	default:
		log.Fatalf("unexpected error: %v", err)
		return status.Errorf(codes.Internal, "internal server error")
	}
}

// ToGRPCMessage converts an error to a string message for gRPC responses.
func ToGRPCMessage(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// FromGRPCError converts a gRPC error back to a standard error, extracting the message from the gRPC status.
func FromGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	return errors.New(st.Message())
}
