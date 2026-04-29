// Package main is the entrypoint of the broker application.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/sceredi/co-type/broker/internal/config"
	cfg_utils "github.com/sceredi/co-type/common/config"
)

func main() {
	cfg_utils.Setup()
	port, err := strconv.Atoi(os.Getenv("BROKER_PORT"))
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	addr := fmt.Sprintf(":%d", port)
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	err = config.CreateDiscoveryListener(addr, lis)
	if err != nil {
		log.Fatalf("Failed to start discovery listener: %v", err)
	}
	fmt.Println("Hello from Broker")
}
