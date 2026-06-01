// Package domain defines the core domain entities and errors for the broker application.
package domain

import "errors"

var (
	// ErrServerNotFound is returned when a server is not found.
	ErrServerNotFound = errors.New("server not found")
	// ErrServerAlreadyExists is returned when a server already exists.
	ErrServerAlreadyExists = errors.New("server already exists")
	// ErrServerNameInvalid is returned when a server name is invalid.
	ErrServerNameInvalid = errors.New("server name is invalid")
	// ErrServerAddrInvalid is returned when a server address is invalid.
	ErrServerAddrInvalid = errors.New("server address is invalid")
	// ErrNoAvailableServers is returned when there are no available servers.
	ErrNoAvailableServers = errors.New("no available servers")
)

// Server represents a server in the broker application.
type Server struct {
	Name string
	Addr string
	Port int
	Load int
}

// NewServer creates a new Server instance with the given address.
func NewServer(name, addr string, port int) (*Server, error) {
	if name == "" {
		return nil, ErrServerNameInvalid
	}
	if addr == "" {
		return nil, ErrServerAddrInvalid
	}
	return &Server{Name: name, Addr: addr, Port: port, Load: 0}, nil
}
