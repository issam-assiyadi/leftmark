package domain

import "time"

// Kind is the marker vocabulary leftmark recognizes in source comments.
type Kind string

const (
	KindTODO     Kind = "TODO"
	KindFIXME    Kind = "FIXME"
	KindNOTE     Kind = "NOTE"
	KindQUESTION Kind = "QUESTION"
)

// Status is the lifecycle of a tracked item.
type Status string

const (
	StatusOpen  Status = "open"
	StatusDoing Status = "doing"
	StatusDone  Status = "done"
)

// Item is a marker as shown to a user: either an untracked, ephemeral open
// marker (ID is empty, re-detected fresh on every scan, never persisted), or
// a tracked marker (doing, or done-and-kept) whose ID, status, and text are
// joined in from the external metadata store rather than the source line.
type Item struct {
	ID     string
	Kind   Kind
	Status Status
	File   string
	Line   int
	Text   string

	// Located is true for untracked open items (nothing to locate) and for
	// tracked items whose tag was found during the latest scan. It is false
	// only for a tracked item whose ID exists in the store but could not be
	// found anywhere in the tree — e.g. the comment was deleted by hand, or
	// the current branch checkout doesn't contain it. Such items are kept,
	// not dropped: silently losing user-created status is worse than
	// surfacing an item that needs re-triage.
	Located  bool
	LastSeen time.Time // when a tracked item's tag was last found; zero for open items
}

// Filter narrows a set of items by kind and/or status. An empty slice means
// "no restriction" for that dimension.
type Filter struct {
	Kinds    []Kind
	Statuses []Status
}

// Match reports whether the item satisfies the filter.
func (f Filter) Match(item Item) bool {
	if len(f.Kinds) > 0 && !containsKind(f.Kinds, item.Kind) {
		return false
	}
	if len(f.Statuses) > 0 && !containsStatus(f.Statuses, item.Status) {
		return false
	}
	return true
}

func containsKind(kinds []Kind, k Kind) bool {
	for _, candidate := range kinds {
		if candidate == k {
			return true
		}
	}
	return false
}

func containsStatus(statuses []Status, s Status) bool {
	for _, candidate := range statuses {
		if candidate == s {
			return true
		}
	}
	return false
}
