#!/bin/sh

set -eu

golangci_lint_cache="$1"

mkdir -p "$golangci_lint_cache"
GOLANGCI_LINT_CACHE="$golangci_lint_cache" golangci-lint run
