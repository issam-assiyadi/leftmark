#!/bin/sh

set -eu

unformatted="$(gofmt -l "$@")"

if [ -n "$unformatted" ]; then
  echo "Go files are not formatted. Run 'make fmt'."
  printf '%s\n' "$unformatted"
  exit 1
fi
