package domain_test

import (
	"testing"

	"github.com/issam-assiyadi/leftmark/domain"
)

func syntaxFor(t *testing.T, ext string) domain.CommentSyntax {
	t.Helper()
	s, ok := domain.SyntaxForExt(ext)
	if !ok {
		t.Fatalf("no syntax registered for %q", ext)
	}
	return s
}

func TestDetectMarker(t *testing.T) {
	cases := []struct {
		name     string
		ext      string
		line     string
		wantKind domain.Kind
		wantText string
	}{
		{
			name:     "go line comment",
			ext:      ".go",
			line:     "// TODO: fix this",
			wantKind: domain.KindTODO,
			wantText: "fix this",
		},
		{
			name:     "python line comment",
			ext:      ".py",
			line:     "# FIXME: handle empty input",
			wantKind: domain.KindFIXME,
			wantText: "handle empty input",
		},
		{
			name:     "js line comment",
			ext:      ".js",
			line:     "  // NOTE: this depends on load order",
			wantKind: domain.KindNOTE,
			wantText: "this depends on load order",
		},
		{
			name:     "html block comment",
			ext:      ".html",
			line:     "<!-- QUESTION: is this markup still used? -->",
			wantKind: domain.KindQUESTION,
			wantText: "is this markup still used?",
		},
		{
			name:     "css block comment",
			ext:      ".css",
			line:     "/* NOTE: block form */",
			wantKind: domain.KindNOTE,
			wantText: "block form",
		},
		{
			name:     "shell line comment",
			ext:      ".sh",
			line:     "# TODO: pin the tool version",
			wantKind: domain.KindTODO,
			wantText: "pin the tool version",
		},
		{
			name:     "sql line comment",
			ext:      ".sql",
			line:     "-- FIXME: this join is O(n^2)",
			wantKind: domain.KindFIXME,
			wantText: "this join is O(n^2)",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			syntax := syntaxFor(t, tc.ext)

			kind, rest, ok := domain.DetectMarker(tc.line, syntax)
			if !ok {
				t.Fatalf("DetectMarker(%q) = ok=false, want true", tc.line)
			}
			if kind != tc.wantKind {
				t.Errorf("DetectMarker(%q) kind = %q, want %q", tc.line, kind, tc.wantKind)
			}
			if rest != tc.wantText {
				t.Errorf("DetectMarker(%q) rest = %q, want %q", tc.line, rest, tc.wantText)
			}
		})
	}
}
