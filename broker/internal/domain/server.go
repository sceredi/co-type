// Package domain defines the core domain entities and errors for the broker application.
package domain

import "errors"

var (
	// ErrServerNotFound is returned when a server is not found.
	ErrServerNotFound = errors.New("server not found")
	// ErrServerAlreadyExists is returned when a server already exists.
	ErrServerAlreadyExists = errors.New("server already exists")
	// ErrServerAddrInvalid is returned when a server address is invalid.
	ErrServerAddrInvalid = errors.New("server address is invalid")
)

// Server represents a server in the broker application.
type Server struct {
	// Addr is the address of the server.
	Addr string
}

// NewServer creates a new Server instance with the given address.
func NewServer(addr string) (*Server, error) {
	if addr == "" {
		return nil, ErrServerAddrInvalid
	}
	return &Server{Addr: addr}, nil
}

// CreateServerRequest represents a request to create a new server.
type CreateServerRequest struct {
	// Addr is the address of the server to create.
	Addr string
}
