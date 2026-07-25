#!/bin/sh

set -eu

golangci_lint_cache="$1"

mkdir -p "$golangci_lint_cache"
cd src
GOLANGCI_LINT_CACHE="$golangci_lint_cache" golangci-lint run
