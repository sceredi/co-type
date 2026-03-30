package main

import (
	"fmt"

	"github.com/sceredi/co-type/client/internal/example"
	shared "github.com/sceredi/co-type/common/example"
)

func main() {
	fmt.Printf("Hello from %s %s!", example.Name, shared.Name)
}
