// Package main is the entry point of the application.
package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/sceredi/co-type/client/internal/tui"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	m := tui.New()

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
