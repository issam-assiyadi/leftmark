// Package fsrw implements application.FileReadWriter with atomic,
// mode-preserving writes to real source files.
package fsrw

import (
	"fmt"
	"os"
	"path/filepath"
)

type ReadWriter struct{}

func New() ReadWriter { return ReadWriter{} }

func (ReadWriter) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Write atomically replaces path's contents with data: a temp file in the
// same directory (avoids cross-device rename failures) is written, synced,
// and renamed over the original, preserving its permission bits. path must
// already exist - refuses to touch a symlink rather than silently
// replacing it with a regular file.
func (ReadWriter) Write(path string, data []byte) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("fsrw: refusing to rewrite symlink %s", path)
	}

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".leftmark-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	if err := writeSyncClose(tmp, data, info.Mode()); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func writeSyncClose(f *os.File, data []byte, mode os.FileMode) error {
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Chmod(mode); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
