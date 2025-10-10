# mcpc build automation.

MODULE  := github.com/adam-eques/mcpc
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

.PHONY: all build test race cover vet lint bench install clean tidy

all: vet test build

build: ## Build the mcpc binary into ./bin
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/ ./cmd/mcpc

install: ## Install mcpc into GOBIN
	go install -ldflags "$(LDFLAGS)" ./cmd/mcpc

test: ## Run the test suite
	go test ./...

race: ## Run tests with the race detector
	go test -race ./...

cover: ## Produce a coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint if installed
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed; skipping"

bench: ## Run benchmarks
	go test -run '^$$' -bench . -benchmem ./...

tidy: ## Tidy modules
	go mod tidy

clean: ## Remove build artefacts
	rm -rf bin dist coverage.out
