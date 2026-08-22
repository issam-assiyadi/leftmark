package application

// Service is the use-case API every frontend (TUI, CLI/JSON adapter, or any
// future Go frontend) calls through. It has no knowledge of how it's
// presented.
type Service struct {
	root   string
	walker FileWalker
	reader FileReader
}

// NewService wires a Service directly to its ports. Most callers should
// prefer the leftmark.New(root) facade; use this constructor to substitute
// a port (e.g. an in-memory FileReader in tests).
func NewService(root string, walker FileWalker, reader FileReader) *Service {
	return &Service{
		root:   root,
		walker: walker,
		reader: reader,
	}
}
