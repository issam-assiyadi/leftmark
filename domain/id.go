package domain

import (
	"crypto/rand"
)

const idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
const idSuffixLen = 6

// NewID mints a short unique ID not present in existing. Callers scanning a
// whole tree must pass the complete set of IDs already known anywhere in
// that tree, and must add each minted ID to that set immediately, so two
// markers newly discovered in the same scan can't collide with each other.
func NewID(existing map[string]struct{}) string {
	for {
		id := "lm-" + randomSuffix(idSuffixLen)
		if _, taken := existing[id]; !taken {
			return id
		}
	}
}

func randomSuffix(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic("domain: failed to read random bytes: " + err.Error())
	}
	out := make([]byte, n)
	for i, b := range buf {
		out[i] = idAlphabet[int(b)%len(idAlphabet)]
	}
	return string(out)
}
