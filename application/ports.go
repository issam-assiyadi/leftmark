package application

// FileWalker visits every regular file under root that fn should consider.
// Implementations are expected to be .gitignore-aware and to always skip
// .git/ regardless of ignore rules.
type FileWalker interface {
	Walk(root string, fn func(path string) error) error
}

// FileReader reads a single file's contents.
type FileReader interface {
	Read(path string) ([]byte, error)
}
