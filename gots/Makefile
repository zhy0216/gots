.PHONY: build test test-unit test-integration clean run emit-go install help

# Binary name
BINARY := gots

# Build the gots binary
build:
	go build -o $(BINARY) ./cmd/gots

# Run all tests (unit + integration, excluding test/ directory with generated .go files)
test:
	go test -v $(shell go list ./... | grep -v /test)

# Run only unit tests (exclude integration)
test-unit:
	go test -v ./pkg/...

# Run only integration tests
test-integration:
	go test -v -run TestIntegration

# Run a specific integration test (usage: make test-file FILE=example)
test-file:
	go test -v -run "TestIntegration/$(FILE)"

# Run tests with coverage
test-cover:
	go test -cover ./...

# Run tests with detailed coverage report
test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
clean:
	rm -f $(BINARY) coverage.out coverage.html
	go clean -testcache

# Run a .gts file (usage: make run FILE=test/example.gts)
run:
	go run ./cmd/gots run $(FILE)

# Emit Go code for a .gts file (usage: make emit-go FILE=test/example.gts)
emit-go:
	go run ./cmd/gots emit-go $(FILE)

# Build and install to $GOPATH/bin
install:
	go install ./cmd/gots

# Start REPL
repl:
	go run ./cmd/gots repl

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks (fmt, vet, test)
check: fmt vet test

# Show help
help:
	@echo "gots Makefile commands:"
	@echo ""
	@echo "  make build           Build the gots binary"
	@echo "  make test            Run all tests"
	@echo "  make test-unit       Run unit tests only"
	@echo "  make test-integration Run integration tests only"
	@echo "  make test-file FILE=example  Run specific integration test"
	@echo "  make test-cover      Run tests with coverage summary"
	@echo "  make test-cover-html Generate HTML coverage report"
	@echo "  make clean           Remove build artifacts"
	@echo "  make run FILE=path   Run a .gts file"
	@echo "  make emit-go FILE=path  Emit Go code for a .gts file"
	@echo "  make install         Install to GOPATH/bin"
	@echo "  make repl            Start interactive REPL"
	@echo "  make fmt             Format Go code"
	@echo "  make vet             Run go vet"
	@echo "  make check           Run fmt, vet, and tests"
