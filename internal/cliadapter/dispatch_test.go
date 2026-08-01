package cliadapter_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/issam-assiyadi/leftmark/internal/cliadapter"
)

func newGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if out, err := exec.Command("git", "init", "-q", dir).CombinedOutput(); err != nil {
		t.Skipf("git not available, skipping: %v: %s", err, out)
	}
	return dir
}

type item struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Status  string `json:"status"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Text    string `json:"text"`
	Located bool   `json:"located"`
}

func run(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errBuf bytes.Buffer
	handled, exitCode := cliadapter.Dispatch(args, &out, &errBuf)
	if !handled {
		t.Fatalf("Dispatch(%v) not recognized as a subcommand", args)
	}
	return out.String(), errBuf.String(), exitCode
}

func TestUnknownArgsFallThroughToTUI(t *testing.T) {
	var out, errBuf bytes.Buffer
	handled, _ := cliadapter.Dispatch(nil, &out, &errBuf)
	if handled {
		t.Errorf("Dispatch(nil) should not be handled, so the caller falls back to the TUI")
	}
}

func TestScanListPromoteResolveOverCLI(t *testing.T) {
	dir := newGitRepo(t)
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\n\n// TODO: fix this\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	stdout, stderr, code := run(t, "scan", "--root", dir, "--json")
	if code != 0 {
		t.Fatalf("scan exit=%d stderr=%s", code, stderr)
	}
	var scanned []item
	if err := json.Unmarshal([]byte(stdout), &scanned); err != nil {
		t.Fatalf("unmarshal scan output %q: %v", stdout, err)
	}
	if len(scanned) != 1 || scanned[0].Status != "open" || scanned[0].ID != "" {
		t.Fatalf("scan output = %+v, want one open item with no ID", scanned)
	}

	stdout, stderr, code = run(t, "promote", "--root", dir, "--json", scanned[0].File, strconv.Itoa(scanned[0].Line))
	if code != 0 {
		t.Fatalf("promote exit=%d stderr=%s", code, stderr)
	}
	var promoted []item
	if err := json.Unmarshal([]byte(stdout), &promoted); err != nil {
		t.Fatalf("unmarshal promote output %q: %v", stdout, err)
	}
	if len(promoted) != 1 || promoted[0].Status != "doing" || promoted[0].ID == "" {
		t.Fatalf("promote output = %+v, want one doing item with an ID", promoted)
	}
	id := promoted[0].ID

	stdout, stderr, code = run(t, "list", "--root", dir, "--json", "--status=doing")
	if code != 0 {
		t.Fatalf("list exit=%d stderr=%s", code, stderr)
	}
	var listed []item
	if err := json.Unmarshal([]byte(stdout), &listed); err != nil {
		t.Fatalf("unmarshal list output %q: %v", stdout, err)
	}
	if len(listed) != 1 || listed[0].ID != id {
		t.Fatalf("list --status=doing = %+v, want the promoted item", listed)
	}

	stdout, stderr, code = run(t, "resolve", "--root", dir, "--json", "--delete", id)
	if code != 0 {
		t.Fatalf("resolve exit=%d stderr=%s", code, stderr)
	}
	var resolved []item
	if err := json.Unmarshal([]byte(stdout), &resolved); err != nil {
		t.Fatalf("unmarshal resolve output %q: %v", stdout, err)
	}
	if len(resolved) != 1 || resolved[0].Status != "done" || resolved[0].Located {
		t.Fatalf("resolve output = %+v, want one done, unlocated item", resolved)
	}

	got, err := os.ReadFile(mainGo)
	if err != nil {
		t.Fatalf("read after resolve: %v", err)
	}
	if string(got) != "package main\n\n" {
		t.Errorf("file after delete-resolve = %q, want the TODO line gone", got)
	}
}
