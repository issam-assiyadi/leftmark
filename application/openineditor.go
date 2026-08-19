package application

// OpenInEditor opens file:line in the user's editor.
func (s *Service) OpenInEditor(file string, line int) error {
	return s.editor.Open(file, line)
}
