package application

import "github.com/issam-assiyadi/leftmark/domain"

// List filters the in-memory snapshot produced by the last Scan. It does
// not re-walk the tree.
func (s *Service) List(filter domain.Filter) []domain.Item {
	s.mu.Lock()
	items := s.items
	s.mu.Unlock()

	var out []domain.Item
	for _, item := range items {
		if filter.Match(item) {
			out = append(out, item)
		}
	}
	return out
}
