// Package config is responsible for setting up the configuration for the broker application.
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
	"time"

	handler "github.com/sceredi/co-type/broker/internal/grpc"
	"github.com/sceredi/co-type/broker/internal/repository/memory"
	"github.com/sceredi/co-type/broker/internal/service"
	"github.com/sceredi/co-type/common/proto/control"
	"github.com/sceredi/co-type/common/proto/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func createListener(grpcServer *grpc.Server, serviceName string) {
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

// CreateListeners creates and starts the gRPC servers.
func CreateListeners() *grpc.Server {
	serverRepo := memory.NewServerRepository()
	serverSvc := service.NewServerService(serverRepo)

	controlHandler := handler.NewControlHandler(serverSvc)
	discoveryHandler := handler.NewDiscoveryHandler(serverSvc)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    3 * time.Second,
			Timeout: 3 * time.Second,
		}),
	)

	control.RegisterControlServiceServer(grpcServer, controlHandler)
	discovery.RegisterDiscoveryServiceServer(grpcServer, discoveryHandler)

	reflection.Register(grpcServer)

	createListener(grpcServer, "control")
	createListener(grpcServer, "discovery")
	return grpcServer
}
