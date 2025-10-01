.PHONY: test test-unit test-services test-handlers test-coverage test-watch test-specific help fmt lint ci-deps clean

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests (services + handlers)"
	@echo "  test-services - Run service layer tests"
	@echo "  test-handlers - Run handler layer tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-watch    - Run tests in watch mode"
	@echo "  test-specific - Run specific test (TEST=TestName)"
	@echo "  fmt           - Format code with gofmt"
	@echo "  lint          - Run golangci-lint"
	@echo "  ci-deps       - Install CI dependencies"
	@echo "  clean         - Clean test artifacts"

# Run all tests
test:
	go test ./... -v

# Run unit tests (fast, no external dependencies)
test-unit:
	go test ./internal/services ./internal/handlers -v

# Run service layer tests specifically
test-services:
	go test ./internal/services -v

# Run handler tests specifically
test-handlers:
	go test ./internal/handlers -v

# Run tests with coverage report
test-coverage:
	go test ./internal/services ./internal/handlers -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in watch mode (requires gotestsum)
test-watch:
	@which gotestsum > /dev/null || (echo "Installing gotestsum..." && go install github.com/gotesttools/gotestsum@latest)
	gotestsum --watch --format testname ./internal/services ./internal/handlers

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html

# Run specific test (usage: make test-specific TEST=TestCompleteItemPreservesIsRead)
test-specific:
	go test ./... -v -run $(TEST)