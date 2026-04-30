// Package grpc contains all the gateway that the client uses to connect to the broker and game servers.
package grpc

import (
	"errors"

	"google.golang.org/grpc/status"
)

func fromGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	return errors.New(st.Message())
}
