package domain

// Line is a single line of a file, split so its exact terminator survives a
// round trip through SplitLines/JoinLines untouched.
type Line struct {
	Content    string
	Terminator string // "\n", "\r\n", or "" for a final line with no trailing newline
}

// SplitLines splits data into lines, preserving each line's terminator
// (including the CRLF-vs-LF distinction and the presence/absence of a
// trailing newline) so JoinLines(SplitLines(data)) reproduces data exactly.
func SplitLines(data []byte) []Line {
	var lines []Line
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] != '\n' {
			continue
		}
		end := i
		term := "\n"
		if end > start && data[end-1] == '\r' {
			end--
			term = "\r\n"
		}
		lines = append(lines, Line{Content: string(data[start:end]), Terminator: term})
		start = i + 1
	}
	if start < len(data) {
		lines = append(lines, Line{Content: string(data[start:]), Terminator: ""})
	}
	return lines
}

// JoinLines is the inverse of SplitLines.
func JoinLines(lines []Line) []byte {
	var out []byte
	for _, l := range lines {
		out = append(out, l.Content...)
		out = append(out, l.Terminator...)
	}
	return out
}
