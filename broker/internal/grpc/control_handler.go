package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/sceredi/co-type/broker/internal/service"
	"github.com/sceredi/co-type/common/proto/control"
)

// ControlHandler implements the ControlServiceServer interface and handles incoming gRPC requests for server management.
type ControlHandler struct {
	control.UnimplementedControlServiceServer
	serverSvc service.ServerService
}

// NewControlHandler creates a new instance of ControlHandler with the provided ServerService.
func NewControlHandler(serverSvc service.ServerService) *ControlHandler {
	return &ControlHandler{serverSvc: serverSvc}
}

// Manage handles incoming gRPC streams for server management. It listens for messages from the stream and processes them accordingly.
func (h *ControlHandler) Manage(stream control.ControlService_ManageServer) error {
	for {
		env, err := stream.Recv()
		if err != nil {
			slog.ErrorContext(context.Background(), "Failed to receive message from stream",
				slog.String("error", err.Error()),
			)
			return err
		}
		slog.InfoContext(context.Background(), "Received message from stream",
			slog.Any("payload", env.GetPayload()),
		)
		switch msg := env.Payload.(type) {
		case *control.ServerEnvelope_Register:
			h.manageRegisterServerReq(msg.Register, stream)
		default:
			slog.WarnContext(context.Background(), "Received unknown message type",
				slog.Any("payload", env.GetPayload()),
			)
			return toGRPCError(errors.New("unknown message type"))
		}
	}
}

func (h *ControlHandler) manageRegisterServerReq(msg *control.RegisterServer, stream control.ControlService_ManageServer) {
	_, err := h.serverSvc.Create(msg.GetHost(), msg.GetPort())
	if err != nil {
		slog.Error("Failed to register server",
			slog.String("host_name", msg.GetHost()),
			slog.Int("port", int(msg.GetPort())),
			slog.String("error", err.Error()),
		)
	}
	err = stream.Send(
		&control.BrokerEnvelope{
			Payload: &control.BrokerEnvelope_RegisterAck{
				RegisterAck: &control.RegisterServerAck{
					Success: err == nil,
					Message: toGRPCMessage(err),
				},
			},
		},
	)
	if err != nil {
		slog.Error("Failed to send register acknowledgment",
			slog.String("host_name", msg.GetHost()),
			slog.Int("port", int(msg.GetPort())),
			slog.String("error", err.Error()),
		)
	}
}
