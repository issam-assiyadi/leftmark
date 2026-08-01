package domain_test

import (
	"strings"
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

func TestDetectAndTagLine(t *testing.T) {
	cases := []struct {
		name     string
		ext      string
		line     string
		wantKind domain.Kind
		wantText string
		id       string
	}{
		{
			name:     "go line comment",
			ext:      ".go",
			line:     "// TODO: fix this",
			wantKind: domain.KindTODO,
			wantText: "fix this",
			id:       "lm-a1b2c3",
		},
		{
			name:     "python line comment",
			ext:      ".py",
			line:     "# FIXME: handle empty input",
			wantKind: domain.KindFIXME,
			wantText: "handle empty input",
			id:       "lm-d4e5f6",
		},
		{
			name:     "js line comment",
			ext:      ".js",
			line:     "  // NOTE: this depends on load order",
			wantKind: domain.KindNOTE,
			wantText: "this depends on load order",
			id:       "lm-g7h8i9",
		},
		{
			name:     "html block comment",
			ext:      ".html",
			line:     "<!-- QUESTION: is this markup still used? -->",
			wantKind: domain.KindQUESTION,
			wantText: "is this markup still used?",
			id:       "lm-j1k2l3",
		},
		{
			name:     "css block comment",
			ext:      ".css",
			line:     "/* NOTE: block form */",
			wantKind: domain.KindNOTE,
			wantText: "block form",
			id:       "lm-m4n5o6",
		},
		{
			name:     "shell line comment",
			ext:      ".sh",
			line:     "# TODO: pin the tool version",
			wantKind: domain.KindTODO,
			wantText: "pin the tool version",
			id:       "lm-p7q8r9",
		},
		{
			name:     "sql line comment",
			ext:      ".sql",
			line:     "-- FIXME: this join is O(n^2)",
			wantKind: domain.KindFIXME,
			wantText: "this join is O(n^2)",
			id:       "lm-s1t2u3",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			syntax := syntaxFor(t, tc.ext)

			// --- untagged: DetectMarker + TagLine ---
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

			tagged, err := domain.TagLine(tc.line, syntax, tc.id)
			if err != nil {
				t.Fatalf("TagLine(%q) error: %v", tc.line, err)
			}

			// If the original line actually used the block-comment form (its
			// close token appears in it), the tag must land before that
			// closer so it stays inside the comment.
			if syntax.BlockClose != "" && strings.Contains(tc.line, syntax.BlockClose) {
				if !strings.HasSuffix(tagged, syntax.BlockClose) {
					t.Errorf("TagLine(%q) = %q, want it to still end with %q", tc.line, tagged, syntax.BlockClose)
				}
			}

			// --- tagged: ParseTag + ParseTaggedLine recover the same marker ---
			gotID, ok := domain.ParseTag(tagged)
			if !ok || gotID != tc.id {
				t.Errorf("ParseTag(%q) = (%q, %v), want (%q, true)", tagged, gotID, ok, tc.id)
			}

			marker, ok := domain.ParseTaggedLine(tagged, syntax)
			if !ok {
				t.Fatalf("ParseTaggedLine(%q) = ok=false, want true", tagged)
			}
			if marker.ID != tc.id || marker.Kind != tc.wantKind || marker.Text != tc.wantText {
				t.Errorf("ParseTaggedLine(%q) = %+v, want {ID:%q Kind:%q Text:%q}",
					tagged, marker, tc.id, tc.wantKind, tc.wantText)
			}

			// --- idempotency: tagging an already-tagged line must fail, never
			// grow a second tag ---
			if _, err := domain.TagLine(tagged, syntax, "lm-shouldnotmint"); err == nil {
				t.Errorf("TagLine(%q) on an already-tagged line succeeded, want an error", tagged)
			}
			if _, _, ok := domain.DetectMarker(tagged, syntax); !ok {
				t.Errorf("DetectMarker(%q) on an already-tagged line = ok=false, want true (still detectable)", tagged)
			}
		})
	}
}

func TestCycle(t *testing.T) {
	cases := []struct {
		in   domain.Status
		want domain.Status
	}{
		{domain.StatusOpen, domain.StatusDoing},
		{domain.StatusDoing, domain.StatusDone},
		{domain.StatusDone, domain.StatusOpen},
	}
	for _, tc := range cases {
		if got := domain.Cycle(tc.in); got != tc.want {
			t.Errorf("Cycle(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNewIDAvoidsCollisions(t *testing.T) {
	existing := map[string]struct{}{}
	for i := 0; i < 1000; i++ {
		id := domain.NewID(existing)
		if _, taken := existing[id]; taken {
			t.Fatalf("NewID returned a duplicate: %q", id)
		}
		existing[id] = struct{}{}
	}
}
