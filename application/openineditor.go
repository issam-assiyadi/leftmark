package application

// OpenInEditor opens file:line in the user's editor. Takes a location
// rather than an item ID because ephemeral open items - the majority of
// what's tracked - have no ID at all; file:line is the one thing every
// Located item has in common.
func (s *Service) OpenInEditor(file string, line int) error {
	return s.editor.Open(file, line)
}
