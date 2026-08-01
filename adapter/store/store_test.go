package store_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/issam-assiyadi/leftmark/adapter/store"
	"github.com/issam-assiyadi/leftmark/application"
	"github.com/issam-assiyadi/leftmark/domain"
)

func newGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("git", "init", "-q", dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("git not available, skipping store test: %v: %s", err, out)
	}
	return dir
}

func TestLoadOnMissingFileReturnsEmpty(t *testing.T) {
	dir := newGitRepo(t)
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	records, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("Load on a fresh repo = %v, want empty", records)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := newGitRepo(t)
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	want := map[string]application.TrackedRecord{
		"lm-abc123": {
			ID:       "lm-abc123",
			Kind:     domain.KindFIXME,
			Status:   domain.StatusDoing,
			Text:     "handle it",
			File:     "main.go",
			Line:     42,
			LastSeen: time.Now().Truncate(time.Second),
		},
	}

	if err := s.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(got) != 1 || got["lm-abc123"].Text != "handle it" || got["lm-abc123"].Line != 42 {
		t.Errorf("round-tripped records = %+v, want %+v", got, want)
	}
	if !got["lm-abc123"].LastSeen.Equal(want["lm-abc123"].LastSeen) {
		t.Errorf("LastSeen = %v, want %v", got["lm-abc123"].LastSeen, want["lm-abc123"].LastSeen)
	}

	storePath := filepath.Join(dir, ".git", "info", "leftmark", "store.json")
	if _, err := os.Stat(storePath); err != nil {
		t.Errorf("expected store file at %s: %v", storePath, err)
	}
}
