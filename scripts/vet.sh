#!/bin/sh

set -eu

go_cache="$1"

mkdir -p "$go_cache"
cd src
GOCACHE="$go_cache" go vet ./...
