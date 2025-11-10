.PHONY: test test-race run-formatters run-lint clean install-formatters install-tools examples

# Run all tests
test:
	@cd remy && go test ./... -v

# Run tests with race detector
test-race:
	@cd remy && go test ./... -race -v

# Format code
run-formatters:
	@golines --base-formatter=gofumpt -w .

# Run linter (if golangci-lint is installed)
run-lint:
	@cd remy && golangci-lint run ./...

# Clean build artifacts and cache
clean:
	@cd remy && go clean -cache -testcache

# Install formatters
install-formatters:
	go install github.com/segmentio/golines@latest
	go install mvdan.cc/gofumpt@latest

# Install development tools (includes formatters)
install-tools: install-formatters
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all examples
examples:
	@echo "Running basic example..."
	@cd examples/basic && go run .
	@echo "\nRunning bindlogger example..."
	@cd examples/bindlogger && go run .
	@echo "\nRunning dynamiconstructor example..."
	@cd examples/dynamiconstructor && go run .
	@echo "\nRunning guessing_types example..."
	@cd examples/guessing_types && go run .

# Default target
.DEFAULT_GOAL := test
