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

# Download Geekdoc theme (pre-built bundle with compiled assets)
docs-theme:
	@echo "Downloading Geekdoc theme..."
	@mkdir -p docs/themes/hugo-geekdoc/
	@curl -L https://github.com/thegeeklab/hugo-geekdoc/releases/latest/download/hugo-geekdoc.tar.gz | tar -xz -C docs/themes/hugo-geekdoc/ --strip-components=1
	@echo "Theme downloaded successfully!"

# Run Hugo development server (downloads theme if needed)
docs-serve: docs-theme-check
	@cd docs && hugo server

# Build Hugo site (downloads theme if needed)
docs-build: docs-theme-check
	@cd docs && hugo

# Check if theme exists, download if missing
docs-theme-check:
	@if [ ! -d "docs/themes/hugo-geekdoc" ]; then \
		$(MAKE) docs-theme; \
	fi

# Default target
.DEFAULT_GOAL := test
