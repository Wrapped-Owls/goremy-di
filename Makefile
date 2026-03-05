.PHONY: test test-race run-formatters run-lint clean \
	install-formatters install-tools \
	examples \
	docs-build docs-serve docs-theme docs-theme-check \
	benchmark benchmark-generate benchmark-run

# --------------------
# Tests
# --------------------

test:
	$(MAKE) -C remy test

test-race:
	$(MAKE) -C remy test-race

run-lint:
	$(MAKE) -C remy lint

clean:
	$(MAKE) -C remy clean

# --------------------
# Formatting / Tools
# --------------------

run-formatters:
	golines --base-formatter=gofumpt -w .

install-formatters:
	go install github.com/segmentio/golines@latest
	go install mvdan.cc/gofumpt@latest

install-tools: install-formatters
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# --------------------
# Delegation
# --------------------

examples:
	$(MAKE) -C examples test

docs-build:
	$(MAKE) -C docs build

docs-serve:
	$(MAKE) -C docs serve

docs-theme:
	$(MAKE) -C docs theme

docs-theme-check:
	$(MAKE) -C docs theme-check

benchmark:
	$(MAKE) -C benchmark bench

# --------------------
# Defaults
# --------------------

.DEFAULT_GOAL := test
