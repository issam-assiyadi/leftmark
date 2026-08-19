package application

import "github.com/issam-assiyadi/leftmark/domain"

// Summary counts a set of scanned items by kind.
type Summary struct {
	Total  int                 `json:"total"`
	ByKind map[domain.Kind]int `json:"by_kind"`
}

// Report summarizes items, as returned by Scan.
func Report(items []domain.Item) Summary {
	summary := Summary{ByKind: make(map[domain.Kind]int)}
	for _, item := range items {
		summary.Total++
		summary.ByKind[item.Kind]++
	}
	return summary
}
