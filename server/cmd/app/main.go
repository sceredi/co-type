// Package main is the entry point of the server application.
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sceredi/co-type/server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config.Setup()
	serverName := os.Getenv("SERVER_NAME")
	serverAddr := os.Getenv("SERVER_ADDR")
	idx := serverName[len(serverName)-1] - '0'
	port := 50050 + int32(idx)
	serverPort := port
	brokerAddr := os.Getenv("BROKER_ADDR")
	brokerPort := os.Getenv("BROKER_PORT")
	addr := fmt.Sprintf("%s:%s", brokerAddr, brokerPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating gRPC client: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error(fmt.Sprintf("Error closing gRPC connection: %v", err))
		}
	}()
	ds := config.CreateDiscoveryService(conn)
	err = ds.Register(serverName, serverAddr, serverPort)
	if err != nil {
		log.Fatalf("Error registering server: %v", err)
	}
}
