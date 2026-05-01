package config

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

// CreateListener sets up a gRPC server listener for the specified service, using the port defined in the environment variable.
func CreateListener(grpcServer *grpc.Server, serviceName string) {
	lc := net.ListenConfig{}
	port, err := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_PORT", strings.ToUpper(serviceName))))
	if err != nil {
		log.Fatalf("Invalid port number for %s service: %v", serviceName, err)
	}
	addr := fmt.Sprintf(":%d", port)
	lis, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Panicf("Failed to listen port %d for service %s: %v", port, serviceName, err)
	}
	slog.InfoContext(context.Background(), "broker listening",
		slog.String("service", serviceName),
		slog.String("port", addr),
	)
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Panicf("Failed to serve gRPC server on port %d for service %s: %v", port, serviceName, err)
		}
	}()
}
