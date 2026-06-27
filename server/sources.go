// Package server exposes the embedded source files of this module so that
// other packages (e.g. internal/snippet) can use them to generate code
// snippets at runtime without any external file access.
package server

import "embed"

// SourceFiles contains all Go source files of this module, embedded at compile time.
//
//go:embed internal cmd
var SourceFiles embed.FS
