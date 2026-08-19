package fsrw_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/issam-assiyadi/leftmark/adapter/fsrw"
)

var fixtures = []struct {
	name     string
	filename string
	content  string
}{
	{
		name:     "crlf endings",
		filename: "main.go",
		content:  "package main\r\n\r\n// TODO: crlf case\r\nfunc main() {}\r\n",
	},
	{
		name:     "no trailing newline",
		filename: "script.py",
		content:  "import os\n\n# FIXME: no eof newline",
	},
}

func TestReadReturnsExactBytes(t *testing.T) {
	r := fsrw.New()

	for _, fx := range fixtures {
		t.Run(fx.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, fx.filename)
			if err := os.WriteFile(path, []byte(fx.content), 0o644); err != nil {
				t.Fatalf("seed fixture: %v", err)
			}

			got, err := r.Read(path)
			if err != nil {
				t.Fatalf("Read: %v", err)
			}
			if string(got) != fx.content {
				t.Errorf("Read = %q, want %q", got, fx.content)
			}
		})
	}
}
