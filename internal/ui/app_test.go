package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/domain"
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
	if got := app.itemsByCategory[domain.KindTODO]; len(got) != 1 || got[0].Text != "fix this" {
		t.Fatalf("itemsByCategory[TODO] = %+v, want one TODO item", got)
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
	if len(app.itemsByCategory[domain.KindTODO]) != 1 || len(app.itemsByCategory[domain.KindFIXME]) != 1 {
		t.Fatalf("itemsByCategory after rescan = %+v, want one TODO and one FIXME", app.itemsByCategory)
	}
}

func TestCategoryMoveUpDownClampToCategoryBounds(t *testing.T) {
	dir := t.TempDir()

	app := newTestApp(t, dir)

	if err := app.categoryMoveUp(nil, nil); err != nil {
		t.Fatalf("categoryMoveUp: %v", err)
	}
	if app.CategorySelected != 0 {
		t.Errorf("CategorySelected = %d after categoryMoveUp at top, want clamped to 0", app.CategorySelected)
	}

	for i := 0; i < len(categoryOrder)+2; i++ {
		if err := app.categoryMoveDown(nil, nil); err != nil {
			t.Fatalf("categoryMoveDown: %v", err)
		}
	}
	if app.CategorySelected != len(categoryOrder)-1 {
		t.Errorf("CategorySelected = %d after moving past the end, want clamped to %d", app.CategorySelected, len(categoryOrder)-1)
	}
}

func TestItemMoveUpDownClampToItemBounds(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	content := "// TODO: one\n// TODO: two\n// TODO: three\n"
	if err := os.WriteFile(mainGo, []byte(content), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	items := app.currentCategoryItems()
	if len(items) != 3 {
		t.Fatalf("currentCategoryItems = %+v, want 3", items)
	}

	if err := app.itemMoveUp(nil, nil); err != nil {
		t.Fatalf("itemMoveUp: %v", err)
	}
	if app.ItemSelected != 0 {
		t.Errorf("ItemSelected = %d after itemMoveUp at top, want clamped to 0", app.ItemSelected)
	}

	for i := 0; i < 5; i++ {
		if err := app.itemMoveDown(nil, nil); err != nil {
			t.Fatalf("itemMoveDown: %v", err)
		}
	}
	if app.ItemSelected != len(items)-1 {
		t.Errorf("ItemSelected = %d after moving past the end, want clamped to %d", app.ItemSelected, len(items)-1)
	}
}

func TestCategoryMoveResetsItemSelection(t *testing.T) {
	dir := t.TempDir()
	content := "// TODO: one\n// TODO: two\n// FIXME: three\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if err := app.itemMoveDown(nil, nil); err != nil {
		t.Fatalf("itemMoveDown: %v", err)
	}
	if app.ItemSelected != 1 {
		t.Fatalf("ItemSelected = %d before category change, want 1", app.ItemSelected)
	}

	if err := app.categoryMoveDown(nil, nil); err != nil {
		t.Fatalf("categoryMoveDown: %v", err)
	}
	if app.ItemSelected != 0 {
		t.Errorf("ItemSelected = %d after categoryMoveDown, want reset to 0", app.ItemSelected)
	}
}
