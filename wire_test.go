package leftmark_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/domain"
)

func newGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", dir).CombinedOutput(); err != nil {
		t.Skipf("git not available, skipping end-to-end test: %v: %s", err, out)
	}
	return dir
}

// TestEndToEnd exercises the whole real stack (fswalk, fsrw, store) wired
// via leftmark.New, the first point in the codebase where every adapter is
// real rather than a fake or an isolated per-package test.
func TestEndToEnd(t *testing.T) {
	dir := newGitRepo(t)
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\n\n// TODO: fix this\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	svc, err := leftmark.New(dir)
	if err != nil {
		t.Fatalf("leftmark.New: %v", err)
	}

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(items) != 1 || items[0].ID != "" || items[0].Status != domain.StatusOpen {
		t.Fatalf("Scan = %+v, want one ephemeral open item", items)
	}
	open := items[0]

	promoted, err := svc.PromoteToDoing(open.File, open.Line)
	if err != nil {
		t.Fatalf("PromoteToDoing: %v", err)
	}
	if promoted.ID == "" || promoted.Status != domain.StatusDoing {
		t.Fatalf("PromoteToDoing = %+v, want a doing item with an ID", promoted)
	}

	got, err := os.ReadFile(mainGo)
	if err != nil {
		t.Fatalf("read after promote: %v", err)
	}
	if string(got) != "package main\n\n// TODO: fix this ["+promoted.ID+"]\n" {
		t.Errorf("file after promote = %q", got)
	}

	resolved, err := svc.ResolveDone(promoted.ID, false)
	if err != nil {
		t.Fatalf("ResolveDone: %v", err)
	}
	if resolved.Status != domain.StatusDone || resolved.Located {
		t.Fatalf("ResolveDone = %+v, want a done, unlocated (archived) item", resolved)
	}

	got, err = os.ReadFile(mainGo)
	if err != nil {
		t.Fatalf("read after resolve: %v", err)
	}
	if string(got) != "package main\n\n" {
		t.Errorf("file after delete-resolve = %q, want the TODO line gone", got)
	}

	storePath := filepath.Join(dir, ".git", "info", "leftmark", "store.json")
	if _, err := os.Stat(storePath); err != nil {
		t.Errorf("expected a store file at %s: %v", storePath, err)
	}
}
