package content

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// this is supposed to be async, how we can handle the pending state in the ui  ??
func LoadPages(dir string) (Pages, error) {
	var pages Pages

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		page, ok, err := pageFormEntry(path, d)

		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		pages = append(pages, page)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return pages, nil
}

func pageFormEntry(path string, d fs.DirEntry) (Page, bool, error) {
	if d.IsDir() {
		return Page{}, false, nil
	}
	if !strings.HasSuffix(d.Name(), ".md") {
		return Page{}, false, nil
	}
	data, err := os.ReadFile(path)

	if err != nil {
		return Page{}, false, err
	}

	name := strings.TrimSuffix(d.Name(), ".md")
	title := normalizeTitle(d.Name())

	return Page{
		Id:      name,
		Title:   title,
		Content: string(data),
		Path:    path,
	}, true, nil
}

func normalizeTitle(fileBase string) string {
	fileBase = strings.ReplaceAll(fileBase, "-", " ")
	fileBase = strings.ReplaceAll(fileBase, "_", " ")
	if len(fileBase) == 0 {
		return "Untitled"
	}
	fileBase = strings.TrimSuffix(fileBase, ".md")

	return strings.ToUpper(fileBase[:1]) + fileBase[1:]
}
