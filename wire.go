// Package leftmark wires the domain/application core to its default
// filesystem-backed adapters — the one-call entry point for embedding
// leftmark as a Go library. Frontends that need a different adapter (e.g.
// a test double) should call application.NewService directly instead.
package leftmark

import (
	"github.com/issam-assiyadi/leftmark/adapter/editor"
	"github.com/issam-assiyadi/leftmark/adapter/fsrw"
	"github.com/issam-assiyadi/leftmark/adapter/fswalk"
	"github.com/issam-assiyadi/leftmark/application"
)

// New wires a fully-functional Service rooted at root.
func New(root string) *application.Service {
	return application.NewService(root, fswalk.New(), fsrw.New(), editor.New())
}
