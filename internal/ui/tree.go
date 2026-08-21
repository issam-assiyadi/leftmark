package ui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/issam-assiyadi/leftmark/domain"
)

type rowKind int

const (
	rowKindDir rowKind = iota
	rowKindFile
	rowKindItem
)

// row is one visible line of the content pane's tree: a directory, a file,
// or a marker item nested under a file. Its index in the []row slice
// buildRows returns is exactly the line it renders to, since each row is
// one Fprintln and domain.Item.Text can never contain a newline.
type row struct {
	kind     rowKind
	label    string // dir/file: leaf name; unused for item rows
	path     string // dir/file: full relative path, forward-slash-joined - collapsed map key
	prefix   string // ancestor indent/branch glyphs, not including this row's own connector
	isLast   bool   // true if this is the last child among its siblings
	expanded bool   // dir/file rows only: true if not collapsed
	hasKids  bool   // dir/file rows only: true if it has children to show
	item     domain.Item
}

type dirNode struct {
	subdirs map[string]*dirNode
	files   map[string]*fileNode
}

type fileNode struct {
	items []domain.Item
}

func newDirNode() *dirNode {
	return &dirNode{subdirs: make(map[string]*dirNode), files: make(map[string]*fileNode)}
}

// buildRows groups items by directory (derived from Item.File) into a tree,
// then flattens it into an ordered list of visible rows, skipping the
// children of any path present (and true) in collapsed.
func buildRows(items []domain.Item, collapsed map[string]bool) []row {
	root := newDirNode()
	for _, item := range items {
		parts := strings.Split(filepath.ToSlash(item.File), "/")
		dir := root
		for _, part := range parts[:len(parts)-1] {
			next, ok := dir.subdirs[part]
			if !ok {
				next = newDirNode()
				dir.subdirs[part] = next
			}
			dir = next
		}
		fileName := parts[len(parts)-1]
		f, ok := dir.files[fileName]
		if !ok {
			f = &fileNode{}
			dir.files[fileName] = f
		}
		f.items = append(f.items, item)
	}

	var rows []row
	walkDir(root, "", "", collapsed, &rows)
	return rows
}

func walkDir(dir *dirNode, path, prefix string, collapsed map[string]bool, rows *[]row) {
	subdirNames := sortedKeys(dir.subdirs)
	fileNames := sortedKeys(dir.files)
	total := len(subdirNames) + len(fileNames)
	i := 0

	for _, name := range subdirNames {
		isLast := i == total-1
		i++

		label, node := compressChain(dir.subdirs[name], name)
		childPath := joinPath(path, label)
		expanded := !collapsed[childPath]
		*rows = append(*rows, row{
			kind:     rowKindDir,
			label:    label,
			path:     childPath,
			prefix:   prefix,
			isLast:   isLast,
			expanded: expanded,
			hasKids:  true,
		})

		if expanded {
			walkDir(node, childPath, childPrefix(prefix, isLast), collapsed, rows)
		}
	}

	for _, name := range fileNames {
		isLast := i == total-1
		i++

		childPath := joinPath(path, name)
		f := dir.files[name]
		expanded := !collapsed[childPath]
		*rows = append(*rows, row{
			kind:     rowKindFile,
			label:    name,
			path:     childPath,
			prefix:   prefix,
			isLast:   isLast,
			expanded: expanded,
			hasKids:  len(f.items) > 0,
		})

		if expanded {
			itemPrefix := childPrefix(prefix, isLast)
			nItems := len(f.items)
			for j, item := range f.items {
				*rows = append(*rows, row{
					kind:   rowKindItem,
					prefix: itemPrefix,
					isLast: j == nItems-1,
					item:   item,
				})
			}
		}
	}
}

// compressChain collapses a chain of directories that each contain exactly
// one subdirectory and no marked files into a single "a/b/c"-labeled row,
// mirroring how editors like VS Code merge such chains in their file tree.
// It stops as soon as a directory has any files or more/fewer than one
// subdirectory, returning that directory as the node to recurse into next.
func compressChain(dir *dirNode, label string) (string, *dirNode) {
	for len(dir.subdirs) == 1 && len(dir.files) == 0 {
		var childName string
		var child *dirNode
		for name, node := range dir.subdirs {
			childName, child = name, node
		}
		label = label + "/" + childName
		dir = child
	}
	return label, dir
}

func childPrefix(parentPrefix string, parentIsLast bool) string {
	if parentIsLast {
		return parentPrefix + "    "
	}
	return parentPrefix + "│   "
}

func joinPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "/" + name
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func connector(r row) string {
	if r.isLast {
		return "╰── "
	}
	return "├── "
}

func expandGlyph(r row) string {
	if !r.hasKids {
		return "  "
	}
	if r.expanded {
		return "▾ "
	}
	return "▸ "
}

// formatRow renders r as a single display line, truncated to fit
// contentWidth (measured in terminal columns, not bytes).
func formatRow(r row, contentWidth int) string {
	var text string
	switch r.kind {
	case rowKindDir:
		text = r.prefix + connector(r) + expandGlyph(r) + r.label + "/"
	case rowKindFile:
		text = r.prefix + connector(r) + expandGlyph(r) + r.label
	case rowKindItem:
		text = r.prefix + connector(r) + fmt.Sprintf("%d %s", r.item.Line, sanitizeText(r.item.Text))
	}

	return runewidth.Truncate(text, contentWidth, "…")
}

func sanitizeText(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, s)
}
