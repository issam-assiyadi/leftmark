// Package editor implements application.EditorLauncher by shelling out to
// $EDITOR.
package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Launcher struct{}

func New() Launcher { return Launcher{} }

func (Launcher) Open(path string, line int) error {
	name := os.Getenv("EDITOR")
	if name == "" {
		name = "vi"
	}

	cmd := exec.Command(name, argsFor(name, path, line)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// argsFor builds editor-appropriate args for jumping straight to a line,
// keyed off the editor's base command name; anything unrecognized just
// gets the bare path.
func argsFor(editorCmd, path string, line int) []string {
	switch filepath.Base(editorCmd) {
	case "vim", "nvim", "vi":
		return []string{fmt.Sprintf("+%d", line), path}
	case "code", "code-insiders":
		return []string{"-g", fmt.Sprintf("%s:%d", path, line)}
	case "subl", "sublime_text":
		return []string{fmt.Sprintf("%s:%d", path, line)}
	default:
		return []string{path}
	}
}
