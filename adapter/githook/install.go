// Package githook installs leftmark's informational (never-blocking) git
// hook script.
package githook

import (
	"fmt"
	"os"
	"path/filepath"
)

const hookScript = `#!/bin/sh
# Installed by leftmark. Informational only - never blocks commit/push.
leftmark report || true
`

// Install writes the hook script into dir (e.g. a repo's
// scripts/git-hooks, to be pointed at via `git config core.hooksPath`) as
// hookName (e.g. "pre-push").
func Install(dir, hookName string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, hookName)
	if err := os.WriteFile(path, []byte(hookScript), 0o755); err != nil {
		return fmt.Errorf("githook: install %s: %w", hookName, err)
	}
	return nil
}
