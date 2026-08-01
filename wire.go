// Package leftmark wires the domain/application core to its default
// filesystem-backed adapters — the one-call entry point for embedding
// leftmark as a Go library. Frontends that need a different adapter (e.g.
// a test double) should call application.NewService directly instead.
package leftmark

import (
	"github.com/issam-assiyadi/leftmark/adapter/editor"
	"github.com/issam-assiyadi/leftmark/adapter/fsrw"
	"github.com/issam-assiyadi/leftmark/adapter/fswalk"
	"github.com/issam-assiyadi/leftmark/adapter/store"
	"github.com/issam-assiyadi/leftmark/application"
)

// New wires a fully-functional Service rooted at root. root must be inside
// a git working tree, since the tracked-item store lives under its
// .git/info.
func New(root string) (*application.Service, error) {
	st, err := store.New(root)
	if err != nil {
		return nil, err
	}

	return application.NewService(root, fswalk.New(), fsrw.New(), editor.New(), st), nil
}
