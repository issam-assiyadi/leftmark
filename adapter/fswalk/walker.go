// Package fswalk implements application.FileWalker over the real
// filesystem.
package fswalk

import (
	"io/fs"
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Walker walks a directory tree, honoring the root .gitignore (if any) and
// always skipping .git regardless of ignore rules. It only consults a
// root-level .gitignore, not git's full per-directory cascading semantics -
// enough for a personal dev tool without reimplementing git's ignore rules.
type Walker struct{}

func New() Walker { return Walker{} }

func (Walker) Walk(root string, fn func(path string) error) error {
	ignore, _ := gitignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"))

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}

		if d.IsDir() {
			switch {
			case rel == ".":
				return nil
			case d.Name() == ".git":
				return filepath.SkipDir
			case ignore != nil && ignore.MatchesPath(rel):
				return filepath.SkipDir
			default:
				return nil
			}
		}

		if ignore != nil && ignore.MatchesPath(rel) {
			return nil
		}

		return fn(path)
	})
}
