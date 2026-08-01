package domain

// Cycle advances a status through the fixed lifecycle open -> doing -> done,
// wrapping back to open.
func Cycle(s Status) Status {
	switch s {
	case StatusOpen:
		return StatusDoing
	case StatusDoing:
		return StatusDone
	case StatusDone:
		return StatusOpen
	default:
		return StatusOpen
	}
}
