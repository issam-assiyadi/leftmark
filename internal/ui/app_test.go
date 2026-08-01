package ui

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
		t.Skipf("git not available, skipping: %v: %s", err, out)
	}
	return dir
}

func newTestApp(t *testing.T, dir string) *App {
	t.Helper()
	svc, err := leftmark.New(dir)
	if err != nil {
		t.Fatalf("leftmark.New: %v", err)
	}
	app, err := New(svc)
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	return app
}

func TestPromoteAndResolveDeleteFlow(t *testing.T) {
	dir := newGitRepo(t)
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\n\n// TODO: fix this\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if len(app.Items) != 1 || app.Items[0].Status != domain.StatusOpen {
		t.Fatalf("initial Items = %+v, want one open item", app.Items)
	}

	if err := app.promoteSelected(nil, nil); err != nil {
		t.Fatalf("promoteSelected: %v", err)
	}
	if len(app.Items) != 1 || app.Items[0].Status != domain.StatusDoing || app.Items[0].ID == "" {
		t.Fatalf("after promote, Items = %+v, want one doing item with an ID", app.Items)
	}

	if err := app.beginResolve(nil, nil); err != nil {
		t.Fatalf("beginResolve: %v", err)
	}
	if app.pending == nil {
		t.Fatalf("beginResolve did not set a pending confirmation")
	}

	if err := app.confirmResolve(true)(nil, nil); err != nil {
		t.Fatalf("confirmResolve(delete): %v", err)
	}
	if app.pending != nil {
		t.Errorf("pending confirmation should be cleared after confirming")
	}
	if len(app.Items) != 1 || app.Items[0].Status != domain.StatusDone || app.Items[0].Located {
		t.Fatalf("after delete-resolve, Items = %+v, want one done, unlocated item", app.Items)
	}

	got, err := os.ReadFile(mainGo)
	if err != nil {
		t.Fatalf("read after resolve: %v", err)
	}
	if string(got) != "package main\n\n" {
		t.Errorf("file after delete-resolve = %q, want the TODO line gone", got)
	}
}

func TestResolveKeepLeavesTagInPlace(t *testing.T) {
	dir := newGitRepo(t)
	mainGo := filepath.Join(dir, "main.go")
	original := []byte("package main\n\n// FIXME: handle it\n")
	if err := os.WriteFile(mainGo, original, 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if err := app.promoteSelected(nil, nil); err != nil {
		t.Fatalf("promoteSelected: %v", err)
	}
	if err := app.beginResolve(nil, nil); err != nil {
		t.Fatalf("beginResolve: %v", err)
	}
	if err := app.confirmResolve(false)(nil, nil); err != nil {
		t.Fatalf("confirmResolve(keep): %v", err)
	}

	if len(app.Items) != 1 || app.Items[0].Status != domain.StatusDone || !app.Items[0].Located {
		t.Fatalf("after keep-resolve, Items = %+v, want one done, still-located item", app.Items)
	}

	got, err := os.ReadFile(mainGo)
	if err != nil {
		t.Fatalf("read after resolve: %v", err)
	}
	if len(got) == len(original) {
		t.Errorf("expected the tag to have been added by promote and to remain after keep-resolve")
	}
}

func TestCancelResolveClearsPendingWithoutChangingStatus(t *testing.T) {
	dir := newGitRepo(t)
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\n\n// FIXME: handle it\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	app := newTestApp(t, dir)
	if err := app.promoteSelected(nil, nil); err != nil {
		t.Fatalf("promoteSelected: %v", err)
	}
	if err := app.beginResolve(nil, nil); err != nil {
		t.Fatalf("beginResolve: %v", err)
	}
	if err := app.cancelResolve(nil, nil); err != nil {
		t.Fatalf("cancelResolve: %v", err)
	}

	if app.pending != nil {
		t.Errorf("cancelResolve should clear the pending confirmation")
	}
	if app.Items[0].Status != domain.StatusDoing {
		t.Errorf("status after cancel = %q, want still doing", app.Items[0].Status)
	}
}

func TestMoveUpDownClampToItemBounds(t *testing.T) {
	dir := newGitRepo(t)
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
