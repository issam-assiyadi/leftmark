package cliadapter_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/issam-assiyadi/leftmark/internal/cliadapter"
)

type item struct {
	Kind string `json:"kind"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
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

func TestScanOverCLI(t *testing.T) {
	dir := t.TempDir()
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
	if len(scanned) != 1 || scanned[0].Kind != "TODO" || scanned[0].Text != "fix this" {
		t.Fatalf("scan output = %+v, want one TODO item", scanned)
	}
}

func TestReportOverCLI(t *testing.T) {
	dir := t.TempDir()
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("// TODO: a\n// TODO: b\n// FIXME: c\n"), 0o644); err != nil {
		t.Fatalf("seed fixture: %v", err)
	}

	stdout, stderr, code := run(t, "report", "--root", dir, "--json")
	if code != 0 {
		t.Fatalf("report exit=%d stderr=%s", code, stderr)
	}
	var summary struct {
		Total  int            `json:"total"`
		ByKind map[string]int `json:"by_kind"`
	}
	if err := json.Unmarshal([]byte(stdout), &summary); err != nil {
		t.Fatalf("unmarshal report output %q: %v", stdout, err)
	}
	if summary.Total != 3 || summary.ByKind["TODO"] != 2 || summary.ByKind["FIXME"] != 1 {
		t.Fatalf("report output = %+v, want 3 total, 2 TODO, 1 FIXME", summary)
	}
}
