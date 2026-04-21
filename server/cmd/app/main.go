// Package main is the entry point of the server application.
package main

import (
	"fmt"

	shared "github.com/sceredi/co-type/common/example"
)

func main() {
	fmt.Printf("Hello from %s!", shared.Name)
}
