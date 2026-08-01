package application

import (
	"fmt"
	"path/filepath"

	"github.com/issam-assiyadi/leftmark/domain"
)

// ResolveDone marks a tracked item as done. If keepTag is false, its tagged
// comment line is deleted from the source and the record is archived (kept
// in the store for history, with no location). If keepTag is true, the tag
// stays in place and the record's status simply flips to done. Rescans
// afterward, same reasoning as PromoteToDoing.
func (s *Service) ResolveDone(id string, keepTag bool) (domain.Item, error) {
	records, err := s.store.Load()
	if err != nil {
		return domain.Item{}, err
	}

	record, ok := records[id]
	if !ok {
		return domain.Item{}, fmt.Errorf("application: unknown tracked item %s", id)
	}

	if keepTag {
		record.Status = domain.StatusDone
		records[id] = record
		if err := s.store.Save(records); err != nil {
			return domain.Item{}, err
		}
	} else {
		path := filepath.Join(s.root, record.File)
		data, err := s.rw.Read(path)
		if err != nil {
			return domain.Item{}, err
		}

		lines := domain.SplitLines(data)
		idx := record.Line - 1
		if idx < 0 || idx >= len(lines) {
			return domain.Item{}, fmt.Errorf("application: %s no longer has line %d for %s", record.File, record.Line, id)
		}
		if tag, ok := domain.ParseTag(lines[idx].Content); !ok || tag != id {
			return domain.Item{}, fmt.Errorf("application: %s:%d no longer matches tracked item %s", record.File, record.Line, id)
		}

		lines = append(lines[:idx], lines[idx+1:]...)
		if err := s.rw.Write(path, domain.JoinLines(lines)); err != nil {
			return domain.Item{}, err
		}

		record.Status = domain.StatusDone
		record.Archived = true
		record.Line = 0
		records[id] = record
		if err := s.store.Save(records); err != nil {
			return domain.Item{}, err
		}
	}

	items, err := s.Scan()
	if err != nil {
		return domain.Item{}, err
	}
	item, ok := findByID(items, id)
	if !ok {
		return domain.Item{}, fmt.Errorf("application: resolved item %s not found after rescan", id)
	}
	return item, nil
}
