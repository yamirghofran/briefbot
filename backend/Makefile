.PHONY: test test-unit test-services test-handlers test-integration test-all test-coverage test-coverage-integration test-watch test-specific help fmt lint ci-deps clean swagger swagger-init swagger-validate swagger-fmt

# Default target
help:
	@echo "Available targets:"
	@echo "  test                    - Run unit tests only (default)"
	@echo "  test-unit               - Run unit tests (services + handlers)"
	@echo "  test-services           - Run service layer unit tests"
	@echo "  test-handlers           - Run handler layer tests"
	@echo "  test-integration        - Run integration tests only"
	@echo "  test-all                - Run all tests (unit + integration)"
	@echo "  test-coverage           - Run unit tests with coverage report"
	@echo "  test-coverage-integration - Run integration tests with coverage"
	@echo "  test-watch              - Run tests in watch mode"
	@echo "  test-specific           - Run specific test (TEST=TestName)"
	@echo "  swagger                 - Generate swagger documentation"
	@echo "  swagger-init            - Initialize swagger (first time setup)"
	@echo "  swagger-validate        - Validate swagger spec"
	@echo "  swagger-fmt             - Format swagger comments"
	@echo "  fmt                     - Format code with gofmt"
	@echo "  lint                    - Run golangci-lint"
	@echo "  ci-deps                 - Install CI dependencies"
	@echo "  clean                   - Clean test artifacts"

# Run all unit tests (excludes integration tests)
test:
	go test ./... -v

# Run unit tests (fast, no external dependencies)
test-unit:
	go test ./internal/services ./internal/handlers -v

# Run service layer unit tests specifically
test-services:
	go test ./internal/services -v

# Run handler tests specifically
test-handlers:
	go test ./internal/handlers -v

# Run integration tests only (requires external dependencies)
test-integration:
	@echo "Running integration tests..."
	@echo "Note: Requires R2 credentials and ffmpeg installed"
	go test ./internal/services -tags=integration -v

# Run all tests (unit + integration)
test-all:
	@echo "Running all tests (unit + integration)..."
	go test ./... -v
	go test ./internal/services -tags=integration -v

# Run unit tests with coverage report
test-coverage:
	go test ./internal/services ./internal/handlers -coverprofile=coverage.out -v
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests with coverage report
test-coverage-integration:
	go test ./internal/services -tags=integration -coverprofile=coverage-integration.out -v
	go tool cover -html=coverage-integration.out -o coverage-integration.html
	@echo "Integration coverage report generated: coverage-integration.html"

# Run tests in watch mode (requires gotestsum)
test-watch:
	@which gotestsum > /dev/null || (echo "Installing gotestsum..." && go install github.com/gotesttools/gotestsum@latest)
	gotestsum --watch --format testname ./internal/services ./internal/handlers

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html coverage-integration.out coverage-integration.html

# Run specific test (usage: make test-specific TEST=TestCompleteItemPreservesIsRead)
test-specific:
	go test ./... -v -run $(TEST)

# Generate swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Initialize swagger (first time setup)
swagger-init:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Validate swagger spec
swagger-validate:
	swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal --parseDepth 1

# Format swagger comments
swagger-fmt:
	swag fmt -g cmd/server/main.go