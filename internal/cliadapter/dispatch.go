// Package cliadapter is the machine-readable CLI/JSON surface leftmark
// exposes for anything that can't import a Go package (a VS Code
// extension, an editor plugin) to drive as a subprocess. It's a client of
// application.Service, same tier as internal/ui, not a special case of it.
package cliadapter

import (
	"fmt"
	"io"
)

var subcommands = map[string]func(args []string, stdout, stderr io.Writer) int{
	"scan":   runScan,
	"report": runReport,
	"hook":   runHook,
}

// Dispatch routes a known subcommand (args[0]) to its handler and reports
// whether it recognized one at all, so the caller can fall back to
// launching the TUI when it didn't (e.g. bare `leftmark`, or `leftmark`
// with no args).
func Dispatch(args []string, stdout, stderr io.Writer) (handled bool, exitCode int) {
	if len(args) == 0 {
		return false, 0
	}

	run, ok := subcommands[args[0]]
	if !ok {
		return false, 0
	}
	return true, run(args[1:], stdout, stderr)
}

func printf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}
