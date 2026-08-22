package leftmark

import (
	"github.com/issam-assiyadi/leftmark/adapter/fsrw"
	"github.com/issam-assiyadi/leftmark/adapter/fswalk"
	"github.com/issam-assiyadi/leftmark/application"
)

// New wires a fully-functional Service rooted at root.
func New(root string) *application.Service {
	return application.NewService(root, fswalk.New(), fsrw.New())
}
