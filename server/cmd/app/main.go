// Package main is the entry point of the server application.
package main

import (
	"fmt"

	"github.com/sceredi/co-type/server/internal/example"
)

func main() {
	fmt.Printf("Hello from %s!", example.Name)
}
