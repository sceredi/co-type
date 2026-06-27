// Package snippet provides a helper for generating random code snippets
// from the server module's own Go source files.
//
// This file is just here as a way to get code snippets, a non toy implementation would use a service ad hoc.
//
//nolint:all
package snippet

import (
	"io/fs"
	"math/rand/v2"
	"strings"

	server "github.com/sceredi/co-type/server"
)

// Random reads a random .go source file (excluding test files) from the
// embedded server module sources and returns a random window of 10-30 lines.
func Random() string {
	var goFiles []string

	// WalkDir error is intentionally ignored; an empty slice falls back gracefully.
	fs.WalkDir(server.SourceFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			goFiles = append(goFiles, path)
		}

		return nil
	})

	if len(goFiles) == 0 {
		return ""
	}

	file := goFiles[rand.IntN(len(goFiles))]

	content, err := server.SourceFiles.ReadFile(file)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	numLines := 10 + rand.IntN(21) // random window between 10 and 30 lines

	if len(lines) <= numLines {
		return strings.Join(lines, "\n")
	}

	maxStart := len(lines) - numLines
	start := rand.IntN(maxStart + 1)

	return strings.Join(lines[start:start+numLines], "\n")
}
