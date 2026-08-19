## Changes

<!-- What changed and why. One bullet per meaningful change. -->

## Architecture check

<!-- Tick each box, or explain in Insights why it doesn't apply.
     Rule: domain -> nothing; application -> domain only (new capabilities
     are a port in application/ports.go, implemented by an adapter);
     adapter/* -> domain/application only; internal/ui and internal/cliadapter
     depend on application/domain, not adapter/* directly (except the existing
     githook case in cliadapter). wire.go is the only place that wires
     adapters + application together. -->

- [ ] `domain/` added no import of `application`, `adapter/*`, or `internal/*`
- [ ] `application/` added no import of `adapter/*` or `internal/*` — any new capability is a port in `application/ports.go`, implemented by an adapter
- [ ] `adapter/*` added no import of `internal/ui` or `internal/cliadapter`
- [ ] `internal/ui` / `internal/cliadapter` added no new direct `adapter/*` import outside the existing `githook` case
- [ ] N/A — this PR doesn't touch `domain/`, `application/`, or `adapter/*`

## Insights

<!-- Tradeoffs, risks, or context a reviewer needs.
     Flag any change to marker grammar — the `TODO`/`FIXME`/`NOTE`/`QUESTION`
     vocabulary or how comments are detected. -->

## Demo

<!-- Terminal recording or screenshots for TUI-visible changes.
     CI can't verify these visually. N/A if not applicable. -->

## Related

<!-- Closes #123 -->
