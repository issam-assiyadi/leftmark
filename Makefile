.PHONY: fmt fmt-check vet lint check

GOFILES := $(shell find src -type f -name '*.go')
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

check: fmt-check vet lint
