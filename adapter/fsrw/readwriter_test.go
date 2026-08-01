package fsrw_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issam-assiyadi/leftmark/adapter/fsrw"
	"github.com/issam-assiyadi/leftmark/domain"
)

type fixture struct {
	name     string
	filename string
	ext      string
	content  string // fixture content, exactly as it should exist on disk
	mode     os.FileMode
}

var fixtures = []fixture{
	{
		name:     "go with CRLF endings",
		filename: "main.go",
		ext:      ".go",
		content:  "package main\r\n\r\n// TODO: crlf case\r\nfunc main() {}\r\n",
		mode:     0o640,
	},
	{
		name:     "python with no trailing newline",
		filename: "script.py",
		ext:      ".py",
		content:  "import os\n\n# FIXME: no eof newline",
		mode:     0o600,
	},
	{
		name:     "css with a block comment",
		filename: "style.css",
		ext:      ".css",
		content:  "body { color: red; }\n/* NOTE: block form */\n.other { color: blue; }\n",
		mode:     0o644,
	},
}

func TestTagLineIsNonDestructive(t *testing.T) {
	rw := fsrw.New()

	for _, fx := range fixtures {
		t.Run(fx.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, fx.filename)
			if err := os.WriteFile(path, []byte(fx.content), fx.mode); err != nil {
				t.Fatalf("seed fixture: %v", err)
			}

			syntax, ok := domain.SyntaxForExt(fx.ext)
			if !ok {
				t.Fatalf("no syntax registered for %q", fx.ext)
			}

			data, err := rw.Read(path)
			if err != nil {
				t.Fatalf("Read: %v", err)
			}
			lines := domain.SplitLines(data)

			idx := findMarkerLine(t, lines, syntax)
			tagged, err := domain.TagLine(lines[idx].Content, syntax, "lm-testid1")
			if err != nil {
				t.Fatalf("TagLine: %v", err)
			}
			lines[idx].Content = tagged

			if err := rw.Write(path, domain.JoinLines(lines)); err != nil {
				t.Fatalf("Write: %v", err)
			}

			assertMode(t, path, fx.mode)

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("re-read: %v", err)
			}
			withoutTag := strings.Replace(string(got), " "+domain.FormatTag("lm-testid1"), "", 1)
			if withoutTag != fx.content {
				t.Errorf("tagging changed more than the tagged line:\n got (tag stripped): %q\n want (original):    %q", withoutTag, fx.content)
			}

			// --- done-and-delete path: remove that same line entirely ---
			data2, err := rw.Read(path)
			if err != nil {
				t.Fatalf("Read after tag: %v", err)
			}
			lines2 := domain.SplitLines(data2)
			lines2 = append(lines2[:idx], lines2[idx+1:]...)

			if err := rw.Write(path, domain.JoinLines(lines2)); err != nil {
				t.Fatalf("Write after delete: %v", err)
			}

			assertMode(t, path, fx.mode)

			afterDelete, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("re-read after delete: %v", err)
			}
			wantAfterDelete := removeLine(fx.content, idx)
			if string(afterDelete) != wantAfterDelete {
				t.Errorf("deleting the marker line corrupted the rest of the file:\n got:  %q\n want: %q", afterDelete, wantAfterDelete)
			}
		})
	}
}

func findMarkerLine(t *testing.T, lines []domain.Line, syntax domain.CommentSyntax) int {
	t.Helper()
	for i, l := range lines {
		if _, _, ok := domain.DetectMarker(l.Content, syntax); ok {
			return i
		}
	}
	t.Fatalf("no marker line found in fixture")
	return -1
}

func assertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode() != want {
		t.Errorf("file mode = %v, want %v", info.Mode(), want)
	}
}

// removeLine reproduces, on the original fixture content string, exactly
// what deleting domain.SplitLines(content)[idx] and rejoining should
// produce - used as the independent expected value for the delete case.
func removeLine(content string, idx int) string {
	lines := domain.SplitLines([]byte(content))
	lines = append(lines[:idx], lines[idx+1:]...)
	return string(domain.JoinLines(lines))
}
