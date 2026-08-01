package cliadapter

import "github.com/issam-assiyadi/leftmark/domain"

// itemJSON is the stable wire shape a VS Code extension (or anything else
// driving leftmark as a subprocess) parses - field names are part of that
// contract, not an implementation detail.
type itemJSON struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Status  string `json:"status"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Text    string `json:"text"`
	Located bool   `json:"located"`
}

func toItemJSON(item domain.Item) itemJSON {
	return itemJSON{
		ID:      item.ID,
		Kind:    string(item.Kind),
		Status:  string(item.Status),
		File:    item.File,
		Line:    item.Line,
		Text:    item.Text,
		Located: item.Located,
	}
}
