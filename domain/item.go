package domain

// Kind is the marker vocabulary leftmark recognizes in source comments.
type Kind string

const (
	KindTODO     Kind = "TODO"
	KindFIXME    Kind = "FIXME"
	KindNOTE     Kind = "NOTE"
	KindQUESTION Kind = "QUESTION"
)

// Item is a marker comment as found by a scan: ephemeral, re-detected fresh
// on every scan, never persisted between runs.
type Item struct {
	Kind Kind
	File string
	Line int
	Text string
}
