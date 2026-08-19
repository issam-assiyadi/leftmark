package application

import (
	"path/filepath"
	"sort"

	"github.com/issam-assiyadi/leftmark/domain"
)

// Scan walks the tree once and reports every marker comment found, fresh.
func (s *Service) Scan() ([]domain.Item, error) {
	var items []domain.Item

	err := s.walker.Walk(s.root, func(path string) error {
		syntax, ok := domain.SyntaxForExt(filepath.Ext(path))
		if !ok {
			return nil
		}

		data, err := s.reader.Read(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(s.root, path)
		if err != nil {
			rel = path
		}

		for i, line := range domain.SplitLines(data) {
			if kind, text, ok := domain.DetectMarker(line.Content, syntax); ok {
				items = append(items, domain.Item{
					Kind: kind,
					File: rel,
					Line: i + 1,
					Text: text,
				})
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].File != items[j].File {
			return items[i].File < items[j].File
		}
		return items[i].Line < items[j].Line
	})

	return items, nil
}
