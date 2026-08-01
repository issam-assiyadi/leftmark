// Package store implements application.Store as a single JSON file under
// the repo's shared .git/info directory - local-only, personal-scratch
// tracking state that never leaves the clone and is never touched by
// clone/push/pull, mirroring git's own convention for .git/info/exclude.
package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/issam-assiyadi/leftmark/application"
)

type Store struct {
	path string
}

// New resolves the store's path under the repo's shared .git/info
// directory. It uses `git rev-parse --git-common-dir` rather than
// --git-dir so that multiple worktrees of the same repo share one store
// instead of each getting its own.
func New(dir string) (*Store, error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--git-common-dir").Output()
	if err != nil {
		return nil, fmt.Errorf("store: not a git repository (or git not available): %w", err)
	}

	gitCommonDir := strings.TrimSpace(string(out))
	if !filepath.IsAbs(gitCommonDir) {
		gitCommonDir = filepath.Join(dir, gitCommonDir)
	}

	infoDir := filepath.Join(gitCommonDir, "info", "leftmark")
	if err := os.MkdirAll(infoDir, 0o755); err != nil {
		return nil, err
	}

	return &Store{path: filepath.Join(infoDir, "store.json")}, nil
}

func (s *Store) Load() (map[string]application.TrackedRecord, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return map[string]application.TrackedRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return map[string]application.TrackedRecord{}, nil
	}

	var records map[string]application.TrackedRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("store: parsing %s: %w", s.path, err)
	}
	return records, nil
}

// Save atomically overwrites the store file. Unlike adapter/fsrw, the
// target may not exist yet (first-ever save) so there's no prior mode to
// preserve - this is a small, deliberately separate writer rather than a
// forced reuse of fsrw's user-source-file semantics.
func (s *Store) Save(records map[string]application.TrackedRecord) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".leftmark-store-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	if err := writeSyncClose(tmp, data); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func writeSyncClose(f *os.File, data []byte) error {
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
