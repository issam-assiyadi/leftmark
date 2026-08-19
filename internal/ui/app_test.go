package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issam-assiyadi/leftmark"
)

func newTestApp(t *testing.T, dir string) *App {
	t.Helper()
	svc := leftmark.New(dir)
	app, err := New(svc)
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	return app
}

func TestNewScansOnStartup(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\n\n// TODO: fix this\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if len(app.Items) != 1 || app.Items[0].Text != "fix this" {
		t.Fatalf("initial Items = %+v, want one TODO item", app.Items)
	}
}

func TestRescanPicksUpChanges(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("// TODO: one\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if len(app.Items) != 1 {
		t.Fatalf("initial Items = %+v, want 1", app.Items)
	}

	if err := os.WriteFile(mainGo, []byte("// TODO: one\n// FIXME: two\n"), 0o644); err != nil {
		t.Fatalf("update fixture: %v", err)
	}
	if err := app.rescan(nil, nil); err != nil {
		t.Fatalf("rescan: %v", err)
	}
	if len(app.Items) != 2 {
		t.Fatalf("Items after rescan = %+v, want 2", app.Items)
	}
}

func TestMoveUpDownClampToItemBounds(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	content := "// TODO: one\n// TODO: two\n// TODO: three\n"
	if err := os.WriteFile(mainGo, []byte(content), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if len(app.Items) != 3 {
		t.Fatalf("Items = %+v, want 3", app.Items)
	}

	if err := app.moveUp(nil, nil); err != nil {
		t.Fatalf("moveUp: %v", err)
	}
	if app.Selected != 0 {
		t.Errorf("Selected = %d after moveUp at top, want clamped to 0", app.Selected)
	}

	for i := 0; i < 5; i++ {
		if err := app.moveDown(nil, nil); err != nil {
			t.Fatalf("moveDown: %v", err)
		}
	}
	if app.Selected != len(app.Items)-1 {
		t.Errorf("Selected = %d after moving past the end, want clamped to %d", app.Selected, len(app.Items)-1)
	}
}
