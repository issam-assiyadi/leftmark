package application

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/issam-assiyadi/leftmark/domain"
)

// Scan walks the tree once, joins any tagged markers to their stored
// records, and reports every open (ephemeral) marker fresh. It never mints
// an ID and never writes to a source file — the only write is a refreshed
// snapshot of the store itself, for records whose tag was found again (or
// found for the first time via a hand-written tag).
func (s *Service) Scan() ([]domain.Item, error) {
	records, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool, len(records))
	var items []domain.Item
	now := time.Now()

	err = s.walker.Walk(s.root, func(path string) error {
		syntax, ok := domain.SyntaxForExt(filepath.Ext(path))
		if !ok {
			return nil
		}

		data, err := s.rw.Read(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(s.root, path)
		if err != nil {
			rel = path
		}

		for i, line := range domain.SplitLines(data) {
			lineNo := i + 1

			if marker, ok := domain.ParseTaggedLine(line.Content, syntax); ok {
				record, known := records[marker.ID]
				if known && record.Archived {
					// The tag exists but its record was archived (done and
					// deleted) - the user must have manually restored the
					// comment. Treat it as freshly promoted rather than
					// resurrecting stale archived state.
					known = false
				}
				if !known {
					record = TrackedRecord{ID: marker.ID, Status: domain.StatusDoing}
				}
				record.Kind = marker.Kind
				record.Text = marker.Text
				record.File = rel
				record.Line = lineNo
				record.LastSeen = now
				records[marker.ID] = record
				seen[marker.ID] = true

				items = append(items, domain.Item{
					ID:       record.ID,
					Kind:     record.Kind,
					Status:   record.Status,
					File:     record.File,
					Line:     record.Line,
					Text:     record.Text,
					Located:  true,
					LastSeen: record.LastSeen,
				})
				continue
			}

			if kind, text, ok := domain.DetectMarker(line.Content, syntax); ok {
				items = append(items, domain.Item{
					Kind:    kind,
					Status:  domain.StatusOpen,
					File:    rel,
					Line:    lineNo,
					Text:    text,
					Located: true,
				})
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	for id, record := range records {
		if record.Archived {
			items = append(items, domain.Item{
				ID:     record.ID,
				Kind:   record.Kind,
				Status: record.Status,
				File:   record.File,
				Text:   record.Text,
			})
			continue
		}
		if seen[id] {
			continue // already appended above, during the walk
		}
		items = append(items, domain.Item{
			ID:       record.ID,
			Kind:     record.Kind,
			Status:   record.Status,
			File:     record.File,
			Line:     record.Line,
			Text:     record.Text,
			Located:  false,
			LastSeen: record.LastSeen,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].File != items[j].File {
			return items[i].File < items[j].File
		}
		return items[i].Line < items[j].Line
	})

	if err := s.store.Save(records); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.items = items
	s.mu.Unlock()

	return items, nil
}
