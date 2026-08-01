package application

import (
	"time"

	"github.com/issam-assiyadi/leftmark/domain"
)

// FileWalker visits every regular file under root that fn should consider.
// Implementations are expected to be .gitignore-aware and to always skip
// .git/ regardless of ignore rules.
type FileWalker interface {
	Walk(root string, fn func(path string) error) error
}

// FileReadWriter reads and atomically rewrites a single file. Write must
// preserve everything about the file the caller doesn't explicitly change
// (permissions, encoding, line endings of untouched lines).
type FileReadWriter interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
}

// EditorLauncher opens a file at a given line in the user's editor and
// blocks until the editor exits.
type EditorLauncher interface {
	Open(path string, line int) error
}

// TrackedRecord is the persisted metadata for a single tracked (non-open)
// item. Kind/Text are refreshed from the source on every scan when the tag
// is still located; Status/Archived are owned by the store and only change
// through PromoteToDoing/ResolveDone.
type TrackedRecord struct {
	ID       string        `json:"id"`
	Kind     domain.Kind   `json:"kind"`
	Status   domain.Status `json:"status"` // doing | done
	Text     string        `json:"text"`
	File     string        `json:"file"`
	Line     int           `json:"line"`
	Archived bool          `json:"archived"`
	LastSeen time.Time     `json:"last_seen"`
}

// Store loads and saves the full set of tracked records. Coarse-grained by
// design: the expected scale is one small per-repo file, and atomicity is
// entirely the adapter's responsibility.
type Store interface {
	Load() (map[string]TrackedRecord, error)
	Save(map[string]TrackedRecord) error
}
