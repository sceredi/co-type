// Package config is responsible for setting up the configuration for the server application.
package config

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

// Config holds the configuration settings for the server application.
type Config struct {
	Name        string
	Addr        string
	Port        int
	ControlAddr string
	ControlPort int
}

var (
	globalCfg *Config
	once      sync.Once
)

// Get returns the global configuration instance, loading it if it hasn't been loaded yet.
func Get() *Config {
	once.Do(load)
	return globalCfg
}

func load() {
	serverName := os.Getenv("SERVER_NAME")
	serverAddr := os.Getenv("SERVER_ADDR")
	serverPortStr := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		log.Fatalf("Error parsing server port: %v", err)
	}
	idx := serverName[len(serverName)-1] - '0'
	serverPort = serverPort + int(idx)
	slog.InfoContext(context.Background(), "Server info",
		slog.String("serverName", serverName),
		slog.String("serverAddr", serverAddr),
		slog.Int("idx", int(idx)),
		slog.Int("serverPort", serverPort),
	)
	controlAddr := os.Getenv("CONTROL_ADDR")
	controlPortStr := os.Getenv("CONTROL_PORT")
	controlPort, err := strconv.Atoi(controlPortStr)
	if err != nil {
		log.Fatalf("Error parsing control port: %v", err)
	}
	globalCfg = &Config{
		Name:        serverName,
		Addr:        serverAddr,
		Port:        serverPort,
		ControlAddr: controlAddr,
		ControlPort: controlPort,
	}
}
