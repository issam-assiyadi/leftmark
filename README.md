# leftmark

`leftmark` is a tracker for the TODOs, FIXMEs, notes, and questions you already leave in your codebase.

## Core Idea

Every developer already writes notes while they work — they just write them as comments. A `TODO` scribbled next to a hacky fix, a `FIXME` on code that isn't quite right yet, a `QUESTION` for something to double-check later. That's a real backlog, but it's invisible: scattered across files, easy to forget, impossible to browse or prioritize.

The goal of `leftmark` is not to give you another place to write notes. It's the opposite: stop maintaining a separate notes app, and treat your codebase as the draft it already is. Write your thinking where the code lives. `leftmark` scans it, tracks it, and gives you a dashboard to browse, filter, and manage it — without changing how you write comments.

In short:

- your codebase is the draft
- comments are the capture mechanism, not a separate notes file
- this app is where you review, triage, and track the marks you've left behind

`leftmark` is the scan/track/report engine plus a set of presentation layers on top of it. The first client is a terminal app, but the engine is built so a VS Code extension, an editor plugin, or other front-ends can be added later without changing how items are scanned or tracked.

## How it works

1. **Scan** — `leftmark` walks your project (respecting `.gitignore`) looking for `TODO`, `FIXME`, `NOTE`, and `QUESTION` comments, in any language.
2. **Track** — the first time a marker is found, a short ID and status (`open`) are appended directly into the comment itself, e.g. `// TODO: fix this [lm-a1b2c3 open]`. No external database — the status lives with the code, so it survives edits and refactors.
3. **Manage** — a dashboard lists every tracked item, filterable by kind or status, letting you cycle status (`open` → `doing` → `done`), jump straight to the file and line in your editor, and rescan on demand. The first dashboard is a terminal app; other front-ends can present the same tracked items later.
4. **Report** — an optional, informational-only git hook prints a summary of open items at commit or push time — never blocking, just a nudge.

## Product Direction

- a fixed, opinionated marker vocabulary to start (`TODO`, `FIXME`, `NOTE`, `QUESTION`) — not user-configurable yet
- no separate persistence layer for status — the source file is the single source of truth
- multi-language support via a lightweight comment-syntax table, not a full parser per language
- git hooks are informational by default; blocking/enforcement modes are a later, opt-in step
- editing still happens in the user's own editor — `leftmark` opens `$EDITOR` at the right file:line, it doesn't try to be one
- the scan/track/report engine is kept independent of any single presentation layer, so the terminal app, and any front-end added later, are just clients built on top of it

## Roadmap

The project EPICs live in [docs/epics.md](/Users/iassiyadi/Side-projects/notes.md/docs/epics.md). This document predates the pivot above and needs a rewrite to match — treat it as historical until then.

## Quality Checks

This project uses standard Go quality tooling:

- `gofmt` for formatting
- `go vet` for baseline static analysis
- `golangci-lint` for additional lint checks

Install Go locally, and make sure `golangci-lint` is available in your shell.

Run the checks from the repository root:

```bash
make fmt
make vet
make lint
make check
```

`make check` validates formatting without rewriting files, then runs `go vet` and `golangci-lint`.

GitHub Actions runs the same checks on pull requests and pushes to `main`.

If you want Git to run checks before pushing, configure the repo-managed hooks once:

```bash
git config core.hooksPath scripts/git-hooks
```

The tracked `pre-push` hook runs `make fmt-check` before allowing a push.
