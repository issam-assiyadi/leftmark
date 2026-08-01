package application

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/issam-assiyadi/leftmark/domain"
)

// PromoteToDoing tags the untracked marker at file:line with a freshly
// minted ID, writes the tag into that one line, and records it in the
// store as "doing". file is relative to the service root, matching what
// Scan reports on domain.Item.File. Rescans afterward so the returned item
// and the service's snapshot reflect the tree exactly as it now stands.
func (s *Service) PromoteToDoing(file string, line int) (domain.Item, error) {
	path := filepath.Join(s.root, file)

	syntax, ok := domain.SyntaxForExt(filepath.Ext(path))
	if !ok {
		return domain.Item{}, fmt.Errorf("application: no comment syntax for %s", path)
	}

	data, err := s.rw.Read(path)
	if err != nil {
		return domain.Item{}, err
	}

	lines := domain.SplitLines(data)
	idx := line - 1
	if idx < 0 || idx >= len(lines) {
		return domain.Item{}, fmt.Errorf("application: %s has no line %d", file, line)
	}

	if _, _, ok := domain.DetectMarker(lines[idx].Content, syntax); !ok {
		return domain.Item{}, fmt.Errorf("application: %s:%d is not a marker", file, line)
	}

	records, err := s.store.Load()
	if err != nil {
		return domain.Item{}, err
	}
	existing := make(map[string]struct{}, len(records))
	for id := range records {
		existing[id] = struct{}{}
	}
	id := domain.NewID(existing)

	tagged, err := domain.TagLine(lines[idx].Content, syntax, id)
	if err != nil {
		return domain.Item{}, err
	}
	lines[idx].Content = tagged

	if err := s.rw.Write(path, domain.JoinLines(lines)); err != nil {
		return domain.Item{}, err
	}

	records[id] = TrackedRecord{
		ID:       id,
		Status:   domain.StatusDoing,
		File:     file,
		Line:     line,
		LastSeen: time.Now(),
	}
	if err := s.store.Save(records); err != nil {
		return domain.Item{}, err
	}

	items, err := s.Scan()
	if err != nil {
		return domain.Item{}, err
	}
	item, ok := findByID(items, id)
	if !ok {
		return domain.Item{}, fmt.Errorf("application: promoted item %s not found after rescan", id)
	}
	return item, nil
}
