package application

import "path/filepath"

// ReadFile reads a file given a root-relative path, the form Item.File uses.
func (s *Service) ReadFile(relPath string) ([]byte, error) {
	return s.reader.Read(filepath.Join(s.root, relPath))
}
