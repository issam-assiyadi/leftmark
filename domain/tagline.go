package domain

import (
	"fmt"
	"strings"
)

// TagLine inserts a fresh ID tag into an untagged marker line. For
// block-comment forms the tag is inserted before the closing token, not
// appended after it, so the tag stays inside the comment (otherwise it would
// sit outside the comment and look untagged on every subsequent scan).
func TagLine(line string, syntax CommentSyntax, id string) (string, error) {
	if _, ok := ParseTag(line); ok {
		return "", fmt.Errorf("domain: line already tagged")
	}

	tag := FormatTag(id)

	if syntax.BlockClose != "" {
		if idx := strings.LastIndex(line, syntax.BlockClose); idx >= 0 {
			return line[:idx] + tag + " " + line[idx:], nil
		}
	}

	return line + " " + tag, nil
}

// TaggedMarker is what a single already-tagged comment line reveals on its
// own: its identity, kind, and text. Status and any other lifecycle
// metadata live in the external store and must be joined in by the caller.
type TaggedMarker struct {
	ID   string
	Kind Kind
	Text string
}

// ParseTaggedLine recovers a TaggedMarker from a line that already carries a
// tag. ok is false if the line isn't both tagged and a recognized marker.
func ParseTaggedLine(line string, syntax CommentSyntax) (TaggedMarker, bool) {
	id, ok := ParseTag(line)
	if !ok {
		return TaggedMarker{}, false
	}

	kind, body, ok := DetectMarker(line, syntax)
	if !ok {
		return TaggedMarker{}, false
	}

	text := strings.TrimSpace(strings.Replace(body, FormatTag(id), "", 1))
	return TaggedMarker{ID: id, Kind: kind, Text: text}, true
}
