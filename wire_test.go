package leftmark_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/domain"
)

// TestEndToEnd exercises the whole real stack (fswalk, fsrw) wired via
// leftmark.New, the first point in the codebase where every adapter is real
// rather than a fake or an isolated per-package test.
func TestEndToEnd(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\n\n// TODO: fix this\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	svc := leftmark.New(dir)

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(items) != 1 || items[0].Kind != domain.KindTODO || items[0].Text != "fix this" {
		t.Fatalf("Scan = %+v, want one TODO item", items)
	}
}
