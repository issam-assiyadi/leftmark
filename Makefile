.PHONY: fmt fmt-check vet lint test check-demos check

GOFILES := $(shell find . -type f -name '*.go' -not -path './.cache/*')
CACHE_DIR := $(CURDIR)/.cache
GO_CACHE := $(CACHE_DIR)/go-build
GOLANGCI_LINT_CACHE := $(CACHE_DIR)/golangci-lint

fmt:
	gofmt -w $(GOFILES)

fmt-check:
	./scripts/fmt-check.sh $(GOFILES)

vet:
	./scripts/vet.sh $(GO_CACHE)

lint:
	./scripts/lint.sh $(GOLANGCI_LINT_CACHE)

test:
	./scripts/test.sh $(GO_CACHE)

check-demos:
	./scripts/check-demos.sh

check: fmt-check vet lint test check-demos
