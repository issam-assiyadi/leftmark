#!/bin/sh

set -eu

go_cache="$1"

mkdir -p "$go_cache"
GOCACHE="$go_cache" go test ./...
