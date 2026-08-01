package application_test

import (
	"fmt"
	"testing"
	"time"

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

// fakeFS is an in-memory FileReadWriter.
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

func (f *fakeFS) Write(path string, data []byte) error {
	f.files[path] = data
	return nil
}

// fakeEditor records the last Open call instead of launching anything.
type fakeEditor struct {
	openedFile string
	openedLine int
}

func (e *fakeEditor) Open(path string, line int) error {
	e.openedFile = path
	e.openedLine = line
	return nil
}

// fakeStore is an in-memory Store.
type fakeStore struct {
	records map[string]application.TrackedRecord
}

func newFakeStore() *fakeStore {
	return &fakeStore{records: map[string]application.TrackedRecord{}}
}

func (s *fakeStore) Load() (map[string]application.TrackedRecord, error) {
	out := make(map[string]application.TrackedRecord, len(s.records))
	for k, v := range s.records {
		out[k] = v
	}
	return out, nil
}

func (s *fakeStore) Save(records map[string]application.TrackedRecord) error {
	s.records = records
	return nil
}

func newTestService(fs *fakeFS, store *fakeStore, paths []string) *application.Service {
	return application.NewService("/repo", &fakeWalker{paths: paths}, fs, &fakeEditor{}, store)
}

func TestScanJoinsOpenAndTrackedItems(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("// TODO: fix this\n// FIXME: handle it [lm-abc123]\n")

	store := newFakeStore()
	store.records["lm-abc123"] = application.TrackedRecord{
		ID:     "lm-abc123",
		Kind:   domain.KindFIXME,
		Status: domain.StatusDoing,
	}

	svc := newTestService(fs, store, []string{"/repo/main.go"})

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("Scan returned %d items, want 2: %+v", len(items), items)
	}

	open := items[0]
	if open.ID != "" || open.Status != domain.StatusOpen || open.Line != 1 || open.Text != "fix this" {
		t.Errorf("open item = %+v, want ephemeral open TODO at line 1", open)
	}

	tracked := items[1]
	if tracked.ID != "lm-abc123" || tracked.Status != domain.StatusDoing || !tracked.Located || tracked.Line != 2 {
		t.Errorf("tracked item = %+v, want located doing FIXME at line 2", tracked)
	}
	if tracked.Text != "handle it" {
		t.Errorf("tracked item text = %q, want %q", tracked.Text, "handle it")
	}
}

func TestScanSurfacesOrphanedItemsRatherThanDropThem(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("// nothing tracked here\n")

	store := newFakeStore()
	lastSeen := time.Now().Add(-24 * time.Hour)
	store.records["lm-gone"] = application.TrackedRecord{
		ID:       "lm-gone",
		Kind:     domain.KindTODO,
		Status:   domain.StatusDoing,
		Text:     "used to be here",
		File:     "main.go",
		Line:     7,
		LastSeen: lastSeen,
	}

	svc := newTestService(fs, store, []string{"/repo/main.go"})

	items, err := svc.Scan()
	if err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Scan returned %d items, want 1: %+v", len(items), items)
	}

	orphan := items[0]
	if orphan.ID != "lm-gone" || orphan.Located || orphan.Status != domain.StatusDoing {
		t.Errorf("orphaned item = %+v, want an unlocated doing item, not dropped", orphan)
	}
	if !orphan.LastSeen.Equal(lastSeen) {
		t.Errorf("orphaned item LastSeen = %v, want preserved value %v", orphan.LastSeen, lastSeen)
	}
}

func TestPromoteToDoingTagsAndTracks(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("// TODO: fix this\n")

	store := newFakeStore()
	svc := newTestService(fs, store, []string{"/repo/main.go"})

	if _, err := svc.Scan(); err != nil {
		t.Fatalf("initial Scan error: %v", err)
	}

	item, err := svc.PromoteToDoing("main.go", 1)
	if err != nil {
		t.Fatalf("PromoteToDoing error: %v", err)
	}
	if item.ID == "" || item.Status != domain.StatusDoing || !item.Located {
		t.Errorf("promoted item = %+v, want a located doing item with an ID", item)
	}

	got := string(fs.files["/repo/main.go"])
	want := fmt.Sprintf("// TODO: fix this [%s]\n", item.ID)
	if got != want {
		t.Errorf("file after promote = %q, want %q", got, want)
	}

	if len(store.records) != 1 {
		t.Fatalf("store has %d records, want 1", len(store.records))
	}
}

func TestResolveDoneDeleteRemovesLineAndArchives(t *testing.T) {
	fs := newFakeFS()
	fs.files["/repo/main.go"] = []byte("// keep me\n// FIXME: handle it [lm-abc123]\n// keep me too\n")

	store := newFakeStore()
	store.records["lm-abc123"] = application.TrackedRecord{
		ID: "lm-abc123", Kind: domain.KindFIXME, Status: domain.StatusDoing,
		Text: "handle it", File: "main.go", Line: 2,
	}

	svc := newTestService(fs, store, []string{"/repo/main.go"})
	if _, err := svc.Scan(); err != nil {
		t.Fatalf("initial Scan error: %v", err)
	}

	item, err := svc.ResolveDone("lm-abc123", false)
	if err != nil {
		t.Fatalf("ResolveDone error: %v", err)
	}
	if item.Status != domain.StatusDone || item.Located {
		t.Errorf("resolved item = %+v, want a done, unlocated (archived) item", item)
	}

	got := string(fs.files["/repo/main.go"])
	want := "// keep me\n// keep me too\n"
	if got != want {
		t.Errorf("file after delete-resolve = %q, want %q (surrounding lines untouched)", got, want)
	}
}

func TestResolveDoneKeepTagLeavesSourceUntouched(t *testing.T) {
	fs := newFakeFS()
	original := []byte("// FIXME: handle it [lm-abc123]\n")
	fs.files["/repo/main.go"] = original

	store := newFakeStore()
	store.records["lm-abc123"] = application.TrackedRecord{
		ID: "lm-abc123", Kind: domain.KindFIXME, Status: domain.StatusDoing,
		Text: "handle it", File: "main.go", Line: 1,
	}

	svc := newTestService(fs, store, []string{"/repo/main.go"})
	if _, err := svc.Scan(); err != nil {
		t.Fatalf("initial Scan error: %v", err)
	}

	item, err := svc.ResolveDone("lm-abc123", true)
	if err != nil {
		t.Fatalf("ResolveDone error: %v", err)
	}
	if item.Status != domain.StatusDone || !item.Located {
		t.Errorf("resolved item = %+v, want a done, still-located item", item)
	}

	if string(fs.files["/repo/main.go"]) != string(original) {
		t.Errorf("file after keep-resolve = %q, want untouched %q", fs.files["/repo/main.go"], original)
	}
}
