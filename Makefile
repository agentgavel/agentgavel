# AgentGavel Makefile. Use GOWORK=off when a parent go.work lists unrelated modules.

export GOWORK ?= off

# Local builds inject VERSION (same -X main.version path as .goreleaser.yml).
VERSION ?= $(shell tr -d '[:space:]' < VERSION 2>/dev/null || echo 0.0.0-dev)
LDFLAGS ?= -X main.version=$(VERSION)

.PHONY: build test lint fmt vet help

help:
	@echo "Targets: build test lint fmt vet"
	@echo "VERSION=$(VERSION) (override with VERSION=x.y.z make build)"

build:
	go build -ldflags "$(LDFLAGS)" -o AgentGavel ./cmd/AgentGavel

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.claude/*')

lint: vet
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || go vet ./...

vet:
	go vet ./...
