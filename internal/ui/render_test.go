package ui

import (
	"testing"

	"github.com/issam-assiyadi/leftmark/domain"
)

func TestFormatCategoryRow(t *testing.T) {
	tests := []struct {
		name  string
		width int
		kind  domain.Kind
		count int
		want  string
	}{
		{"fits with room to spare", 20, domain.KindTODO, 3, "TODO              3 "},
		{"exact fit, no truncation", 9, domain.KindTODO, 3, "TODO   3 "},
		{"truncated with ellipsis", 12, domain.Kind("TO CHALLENGE"), 2, "TO CHA... 2 "},
		{"too narrow for ellipsis", 6, domain.Kind("TO CHALLENGE"), 2, "TO  2 "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatCategoryRow(tt.width, tt.kind, tt.count)
			if got != tt.want {
				t.Errorf("formatCategoryRow(%d, %q, %d) = %q, want %q", tt.width, tt.kind, tt.count, got, tt.want)
			}
			if len(got) != tt.width {
				t.Errorf("formatCategoryRow(%d, %q, %d) length = %d, want %d", tt.width, tt.kind, tt.count, len(got), tt.width)
			}
		})
	}
}
