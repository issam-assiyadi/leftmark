package application

import "github.com/issam-assiyadi/leftmark/domain"

// Summary counts the current snapshot by status and kind, and separately
// flags tracked items that couldn't be located in the latest scan.
type Summary struct {
	Total     int
	ByStatus  map[domain.Status]int
	ByKind    map[domain.Kind]int
	Unlocated int
}

func (s *Service) Report() Summary {
	s.mu.Lock()
	items := s.items
	s.mu.Unlock()

	summary := Summary{
		ByStatus: make(map[domain.Status]int),
		ByKind:   make(map[domain.Kind]int),
	}
	for _, item := range items {
		summary.Total++
		summary.ByStatus[item.Status]++
		summary.ByKind[item.Kind]++
		if item.ID != "" && !item.Located {
			summary.Unlocated++
		}
	}
	return summary
}
