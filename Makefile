
.PHONY: help watch run status interfaces test race vet fmt build clean

BINARY := bin/flux

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  make %-12s %s\n", $$1, $$2}'

watch: ## Monitor network traffic
	go run ./cmd/flux watch

run: watch ## Alias for watch

status: ## Check network counter availability
	go run ./cmd/flux status

interfaces: ## List network interfaces
	go run ./cmd/flux interfaces

test: ## Run all tests
	go test ./...

race: ## Run tests with race detection
	go test -race ./...

vet: ## Run Go static analysis
	go vet ./...

fmt: ## Format Go files
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

build: ## Build Flux binary
	mkdir -p bin
	go build -o $(BINARY) ./cmd/flux

clean: ## Remove the built binary
	rm -rf bin
