# AgentGavel Makefile. Use GOWORK=off when a parent go.work lists unrelated modules.

export GOWORK ?= off

.PHONY: build test lint fmt vet help

help:
	@echo "Targets: build test lint fmt vet"

build:
	go build -o AgentGavel ./cmd/AgentGavel

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.claude/*')

lint: vet
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || go vet ./...

vet:
	go vet ./...
