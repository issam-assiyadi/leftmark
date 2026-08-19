package domain

import (
	"strings"
)

var markerKinds = []Kind{KindTODO, KindFIXME, KindNOTE, KindQUESTION}

// DetectMarker reports whether line contains a comment (per syntax) whose
// body starts with one of the recognized marker keywords. rest is everything
// in the comment after the keyword and its separating colon, trimmed, with
// any trailing block-comment closer removed so it never leaks into the text
// shown to a user.
func DetectMarker(line string, syntax CommentSyntax) (kind Kind, rest string, ok bool) {
	start, found := commentBodyStart(line, syntax)
	if !found {
		return "", "", false
	}
	body := strings.TrimSpace(line[start:])
	for _, k := range markerKinds {
		kw := string(k)
		if !strings.HasPrefix(body, kw) {
			continue
		}
		after := body[len(kw):]
		if after != "" && after[0] != ':' && !isSpace(after[0]) {
			continue // e.g. "TODOIST" is not the keyword "TODO"
		}
		after = strings.TrimSpace(strings.TrimPrefix(after, ":"))
		if syntax.BlockClose != "" {
			after = strings.TrimSpace(strings.TrimSuffix(after, syntax.BlockClose))
		}
		return k, after, true
	}
	return "", "", false
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

// commentBodyStart finds the earliest comment opener (line-prefix or block
// form) on the line and returns the index right after it. On ties it prefers
// the longer token, so a block opener that itself begins with a shorter line
// prefix (e.g. Lua's "--[[" vs "--") wins correctly.
func commentBodyStart(line string, syntax CommentSyntax) (int, bool) {
	best := -1
	bestLen := 0

	consider := func(token string) {
		if token == "" {
			return
		}
		i := strings.Index(line, token)
		if i < 0 {
			return
		}
		if best == -1 || i < best || (i == best && len(token) > bestLen) {
			best = i
			bestLen = len(token)
		}
	}

	for _, p := range syntax.LinePrefixes {
		consider(p)
	}
	consider(syntax.BlockOpen)

	if best == -1 {
		return 0, false
	}
	return best + bestLen, true
}
