package cliadapter

import "github.com/issam-assiyadi/leftmark/domain"

// itemJSON is the wire shape a VS Code extension (or anything else driving
// leftmark as a subprocess) parses - field names are part of that contract,
// not an implementation detail.
type itemJSON struct {
	Kind string `json:"kind"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

func toItemJSON(item domain.Item) itemJSON {
	return itemJSON{
		Kind: string(item.Kind),
		File: item.File,
		Line: item.Line,
		Text: item.Text,
	}
}
