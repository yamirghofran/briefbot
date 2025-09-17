.PHONY: test test-unit test-services test-coverage test-watch help

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-services - Run service layer tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-watch    - Run tests in watch mode"

# Run all tests
test:
	go test ./... -v

# Run unit tests (fast, no external dependencies)
test-unit:
	go test ./internal/services -v

# Run service layer tests specifically
test-services:
	go test ./internal/services -v

# Run tests with coverage report
test-coverage:
	go test ./internal/services -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in watch mode (requires gotestsum)
test-watch:
	@which gotestsum > /dev/null || (echo "Installing gotestsum..." && go install github.com/gotesttools/gotestsum@latest)
	gotestsum --watch --format testname ./internal/services

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html

# Run specific test (usage: make test-specific TEST=TestCompleteItemPreservesIsRead)
test-specific:
	go test ./internal/services -v -run $(TEST)