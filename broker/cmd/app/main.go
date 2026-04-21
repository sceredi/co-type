// Package main is the entrypoint of the broker application.
package main

import (
	"fmt"

	"github.com/sceredi/co-type/broker/internal/example"
)

func main() {
	fmt.Printf("Hello from %s!", example.Name)
}
