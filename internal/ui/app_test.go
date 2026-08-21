package ui

import (
	"errors"
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

func TestRowMoveUpDownClampToRowBounds(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	content := "// TODO: one\n// TODO: two\n// TODO: three\n"
	if err := os.WriteFile(mainGo, []byte(content), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	rows := app.visibleRows()
	if len(rows) != 4 {
		t.Fatalf("visibleRows = %+v, want 4 (1 file header + 3 items)", rows)
	}

	if err := app.rowMoveUp(nil, nil); err != nil {
		t.Fatalf("rowMoveUp: %v", err)
	}
	if app.RowSelected != 0 {
		t.Errorf("RowSelected = %d after rowMoveUp at top, want clamped to 0", app.RowSelected)
	}

	for i := 0; i < 5; i++ {
		if err := app.rowMoveDown(nil, nil); err != nil {
			t.Fatalf("rowMoveDown: %v", err)
		}
	}
	if app.RowSelected != len(rows)-1 {
		t.Errorf("RowSelected = %d after moving past the end, want clamped to %d", app.RowSelected, len(rows)-1)
	}
}

func TestCategoryMoveResetsItemSelection(t *testing.T) {
	dir := t.TempDir()
	content := "// TODO: one\n// TODO: two\n// FIXME: three\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if err := app.rowMoveDown(nil, nil); err != nil {
		t.Fatalf("rowMoveDown: %v", err)
	}
	if app.RowSelected != 1 {
		t.Fatalf("RowSelected = %d before category change, want 1", app.RowSelected)
	}

	if err := app.categoryMoveDown(nil, nil); err != nil {
		t.Fatalf("categoryMoveDown: %v", err)
	}
	if app.RowSelected != 0 {
		t.Errorf("RowSelected = %d after categoryMoveDown, want reset to 0", app.RowSelected)
	}
}

func TestCollapseSurvivesRescan(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "a.go"), []byte("// TODO: one\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	app.collapsed["sub"] = true

	before := len(app.visibleRows())
	if err := app.rescan(nil, nil); err != nil {
		t.Fatalf("rescan: %v", err)
	}
	after := len(app.visibleRows())

	if before != after {
		t.Fatalf("visibleRows changed across rescan (%d -> %d), want collapse state preserved", before, after)
	}
	if !app.collapsed["sub"] {
		t.Errorf("collapsed[\"sub\"] = false after rescan, want still collapsed")
	}
}

func TestCollapseSurvivesCategorySwitch(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	content := "// TODO: one\n// FIXME: two\n"
	if err := os.WriteFile(filepath.Join(subdir, "a.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	app.collapsed["sub"] = true

	if err := app.categoryMoveDown(nil, nil); err != nil {
		t.Fatalf("categoryMoveDown: %v", err)
	}
	if err := app.categoryMoveUp(nil, nil); err != nil {
		t.Fatalf("categoryMoveUp: %v", err)
	}

	if !app.collapsed["sub"] {
		t.Errorf("collapsed[\"sub\"] = false after switching category and back, want still collapsed")
	}
}

func TestActivateSelectedRowDirTogglesCollapse(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "sub")
	if err := os.Mkdir(subdir, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subdir, "a.go"), []byte("// TODO: one\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	rows := app.visibleRows()
	if len(rows) == 0 || rows[0].kind != rowKindDir {
		t.Fatalf("visibleRows[0] = %+v, want a dir row", rows)
	}
	app.RowSelected = 0

	if err := app.activateSelectedRow(nil, nil); err != nil {
		t.Fatalf("activateSelectedRow: %v", err)
	}
	if !app.collapsed[rows[0].path] {
		t.Errorf("collapsed[%q] = false after activateSelectedRow on a dir row, want true", rows[0].path)
	}
	if app.pendingEdit != nil {
		t.Errorf("pendingEdit = %+v after toggling a dir row, want nil", app.pendingEdit)
	}
}

func TestActivateSelectedRowItemOpensEditor(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("// TODO: one\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	rows := app.visibleRows()
	itemIdx := -1
	for i, r := range rows {
		if r.kind == rowKindItem {
			itemIdx = i
			break
		}
	}
	if itemIdx < 0 {
		t.Fatalf("visibleRows = %+v, want at least one item row", rows)
	}
	app.RowSelected = itemIdx

	err := app.activateSelectedRow(nil, nil)
	if !errors.Is(err, errOpenEditor) {
		t.Fatalf("activateSelectedRow on item row error = %v, want errOpenEditor", err)
	}
	if app.pendingEdit == nil {
		t.Fatalf("pendingEdit = nil after activateSelectedRow on item row, want set")
	}
}
