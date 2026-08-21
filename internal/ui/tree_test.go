package ui

import (
	"testing"

	"github.com/mattn/go-runewidth"

	"github.com/issam-assiyadi/leftmark/domain"
)

func TestBuildRowsMultiDir(t *testing.T) {
	items := []domain.Item{
		{Kind: domain.KindTODO, File: "b.go", Line: 1, Text: "root file"},
		{Kind: domain.KindTODO, File: "zdir/a.go", Line: 1, Text: "in zdir"},
		{Kind: domain.KindTODO, File: "adir/nested/a.go", Line: 1, Text: "deep"},
		{Kind: domain.KindTODO, File: "adir/a.go", Line: 1, Text: "in adir"},
	}

	rows := buildRows(items, nil)

	type summary struct {
		kind  rowKind
		label string
	}
	var gotRows []summary
	for _, r := range rows {
		if r.kind == rowKindItem {
			continue
		}
		gotRows = append(gotRows, summary{r.kind, r.label})
	}

	want := []summary{
		{rowKindDir, "adir"},
		{rowKindDir, "nested"},
		{rowKindFile, "a.go"}, // adir/nested/a.go
		{rowKindFile, "a.go"}, // adir/a.go
		{rowKindDir, "zdir"},
		{rowKindFile, "a.go"}, // zdir/a.go
		{rowKindFile, "b.go"},
	}

	if len(gotRows) != len(want) {
		t.Fatalf("buildRows dir/file rows = %+v, want %+v", gotRows, want)
	}
	for i := range want {
		if gotRows[i] != want[i] {
			t.Errorf("row %d = %+v, want %+v", i, gotRows[i], want[i])
		}
	}
}

func TestBuildRowsCollapse(t *testing.T) {
	items := []domain.Item{
		{Kind: domain.KindTODO, File: "dir/a.go", Line: 1, Text: "one"},
		{Kind: domain.KindTODO, File: "dir/a.go", Line: 2, Text: "two"},
		{Kind: domain.KindTODO, File: "other.go", Line: 1, Text: "three"},
	}

	expanded := buildRows(items, map[string]bool{})
	// dir, dir/a.go, item one, item two, other.go, other.go's item = 6 rows
	if len(expanded) != 6 {
		t.Fatalf("expanded rows = %d, want 6: %+v", len(expanded), expanded)
	}

	collapsedDir := buildRows(items, map[string]bool{"dir": true})
	// dir, other.go, other.go's item = 3 rows (dir/a.go + its items hidden)
	if len(collapsedDir) != 3 {
		t.Fatalf("collapsed-dir rows = %d, want 3: %+v", len(collapsedDir), collapsedDir)
	}

	collapsedFile := buildRows(items, map[string]bool{"dir/a.go": true})
	// dir, dir/a.go, other.go, other.go's item = 4 rows (dir/a.go's items hidden)
	if len(collapsedFile) != 4 {
		t.Fatalf("collapsed-file rows = %d, want 4: %+v", len(collapsedFile), collapsedFile)
	}
}

func TestBuildRowsCompressesSingleChildChains(t *testing.T) {
	items := []domain.Item{
		{Kind: domain.KindTODO, File: "folder-1/folder-2/sub-folder-1/test.go", Line: 1, Text: "deep"},
	}

	rows := buildRows(items, nil)
	if len(rows) != 3 {
		t.Fatalf("buildRows = %+v, want 3 (1 merged dir + 1 file + 1 item)", rows)
	}

	dirRow := rows[0]
	if dirRow.kind != rowKindDir || dirRow.label != "folder-1/folder-2/sub-folder-1" {
		t.Errorf("rows[0] = %+v, want a dir row labeled \"folder-1/folder-2/sub-folder-1\"", dirRow)
	}
	if dirRow.path != "folder-1/folder-2/sub-folder-1" {
		t.Errorf("rows[0].path = %q, want %q", dirRow.path, "folder-1/folder-2/sub-folder-1")
	}

	fileRow := rows[1]
	if fileRow.kind != rowKindFile || fileRow.label != "test.go" {
		t.Errorf("rows[1] = %+v, want a file row labeled \"test.go\"", fileRow)
	}

	collapsed := buildRows(items, map[string]bool{"folder-1/folder-2/sub-folder-1": true})
	if len(collapsed) != 1 {
		t.Fatalf("collapsing the merged chain = %+v, want just the 1 dir row", collapsed)
	}
}

func TestBuildRowsStopsCompressionAtBranchOrFile(t *testing.T) {
	items := []domain.Item{
		{Kind: domain.KindTODO, File: "a/b/one.go", Line: 1, Text: "one"},
		{Kind: domain.KindTODO, File: "a/b/c/two.go", Line: 1, Text: "two"},
		{Kind: domain.KindTODO, File: "x/y/three.go", Line: 1, Text: "three"},
	}

	rows := buildRows(items, nil)

	var labels []string
	for _, r := range rows {
		if r.kind == rowKindDir {
			labels = append(labels, r.label)
		}
	}

	// "a" has no files of its own and exactly one subdirectory ("b"), so it
	// merges into "a/b". "a/b" itself has a file (one.go) alongside its
	// subdir "c", so the chain stops there - "c" gets its own row instead
	// of merging further. "x/y" has no files at either level, so it merges
	// fully into one row.
	want := []string{"a/b", "c", "x/y"}
	if len(labels) != len(want) {
		t.Fatalf("dir labels = %v, want %v", labels, want)
	}
	for i := range want {
		if labels[i] != want[i] {
			t.Errorf("dir label %d = %q, want %q (all: %v)", i, labels[i], want[i], labels)
		}
	}
}

func TestFormatRowWidth(t *testing.T) {
	tests := []struct {
		name string
		r    row
		w    int
	}{
		{
			name: "deep prefix with box-drawing glyphs",
			r: row{
				kind:    rowKindFile,
				label:   "file.go",
				prefix:  "│   │   ",
				isLast:  true,
				hasKids: true,
			},
			w: 12,
		},
		{
			name: "item row with wide CJK text",
			r: row{
				kind:   rowKindItem,
				prefix: "    ",
				isLast: true,
				item:   domain.Item{Line: 42, Text: "修复这个问题 fix this"},
			},
			w: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatRow(tt.r, tt.w)
			if w := runewidth.StringWidth(got); w > tt.w {
				t.Errorf("formatRow(%+v, %d) = %q, display width %d > %d", tt.r, tt.w, got, w, tt.w)
			}
		})
	}
}

func TestFormatRowWidthWouldFailNaiveLenCheck(t *testing.T) {
	r := row{
		kind:    rowKindDir,
		label:   "dir",
		prefix:  "│   │   │   ",
		isLast:  false,
		hasKids: true,
	}
	w := 20

	got := formatRow(r, w)
	if len(got) <= w {
		t.Fatalf("expected byte length to exceed display width for this glyph-heavy row, got byte len %d for %q", len(got), got)
	}
	if runewidth.StringWidth(got) > w {
		t.Fatalf("formatRow(%+v, %d) = %q, display width %d > %d", r, w, got, runewidth.StringWidth(got), w)
	}
}
