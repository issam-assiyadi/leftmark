package application

import (
	"sync"

	"github.com/issam-assiyadi/leftmark/domain"
)

// Service is the use-case API every frontend (TUI, CLI/JSON adapter, or any
// future Go frontend) calls through. It has no knowledge of how it's
// presented.
type Service struct {
	root   string
	walker FileWalker
	rw     FileReadWriter
	editor EditorLauncher
	store  Store

	mu    sync.Mutex
	items []domain.Item
}

// NewService wires a Service directly to its ports. Most callers should
// prefer the leftmark.New(root) facade; use this constructor to substitute
// a port (e.g. an in-memory Store in tests).
func NewService(root string, walker FileWalker, rw FileReadWriter, editor EditorLauncher, store Store) *Service {
	return &Service{
		root:   root,
		walker: walker,
		rw:     rw,
		editor: editor,
		store:  store,
	}
}

func findByID(items []domain.Item, id string) (domain.Item, bool) {
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return domain.Item{}, false
}
