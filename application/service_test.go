package application_test

import (
	"fmt"
	"testing"

	"github.com/issam-assiyadi/leftmark/application"
	"github.com/issam-assiyadi/leftmark/domain"
)

// fakeWalker walks a fixed, in-memory list of paths.
type fakeWalker struct {
	paths []string
}

func (w *fakeWalker) Walk(root string, fn func(path string) error) error {
	for _, p := range w.paths {
		if err := fn(p); err != nil {
			return err
		}
	}
	return nil
}

// fakeFS is an in-memory FileReader.
type fakeFS struct {
	files map[string][]byte
}

func newFakeFS() *fakeFS { return &fakeFS{files: map[string][]byte{}} }

func (f *fakeFS) Read(path string) ([]byte, error) {
	data, ok := f.files[path]
	if !ok {
		return nil, fmt.Errorf("fakeFS: no such file %s", path)
	}
	return data, nil
}

func newTestService(fs *fakeFS, paths []string) *application.Service {
	return application.NewService("/repo", &fakeWalker{paths: paths}, fs)
}

func TestScanDetectsMarkersAcrossKinds(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("// TODO: fix this\n// FIXME: handle it\n// NOTE: heads up\n")

	svc := newTestService(fs, []string{"/repo/main.go"})

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("Scan returned %d items, want 3: %+v", len(items), items)
	}

	if items[0].Kind != domain.KindTODO || items[0].Line != 1 || items[0].Text != "fix this" {
		t.Errorf("items[0] = %+v, want TODO at line 1", items[0])
	}
	if items[1].Kind != domain.KindFIXME || items[1].Line != 2 || items[1].Text != "handle it" {
		t.Errorf("items[1] = %+v, want FIXME at line 2", items[1])
	}
	if items[2].Kind != domain.KindNOTE || items[2].Line != 3 || items[2].Text != "heads up" {
		t.Errorf("items[2] = %+v, want NOTE at line 3", items[2])
	}
}

func TestScanIgnoresUnrecognizedFilesAndLines(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("package main\n\nfunc main() {}\n")

	svc := newTestService(fs, []string{"/repo/main.go"})

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("Scan returned %d items, want 0: %+v", len(items), items)
	}
}

func TestScanSortsByFileThenLine(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/b.go"] = []byte("// TODO: in b\n")
	fs.files["/repo/a.go"] = []byte("// TODO: second in a\n// TODO: first in a\n")

	svc := newTestService(fs, []string{"/repo/b.go", "/repo/a.go"})

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("Scan returned %d items, want 3: %+v", len(items), items)
	}
	if items[0].File != "a.go" || items[0].Line != 1 {
		t.Errorf("items[0] = %+v, want a.go line 1 first", items[0])
	}
	if items[1].File != "a.go" || items[1].Line != 2 {
		t.Errorf("items[1] = %+v, want a.go line 2 second", items[1])
	}
	if items[2].File != "b.go" {
		t.Errorf("items[2] = %+v, want b.go last", items[2])
	}
}

func TestReadFileJoinsRoot(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("package main\n")
	svc := newTestService(fs, nil)

	data, err := svc.ReadFile("main.go")
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if string(data) != "package main\n" {
		t.Errorf("ReadFile = %q, want file contents", data)
	}
}
