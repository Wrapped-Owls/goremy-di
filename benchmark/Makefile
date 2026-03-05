.PHONY: install-tools generate benchstat-check \
	bench-files bench-remy bench-do bench-dig \
	benchstat update-readme bench clean

BENCH_DIR := .benchstat
BENCH_FLAGS := -run=^$$ -bench=. -benchmem -count=5 ./tests

DIG_FILE := dig.txt
DO_FILE := do.txt
REMY_FILE := remy.txt

install-tools:
	go install golang.org/x/perf/cmd/benchstat@latest

generate:
	go generate ./...

benchstat-check:
	@command -v benchstat >/dev/null 2>&1 || (echo "benchstat is required. Run: make install-tools" && exit 1)

bench-files: generate bench-remy bench-do bench-dig

bench-remy:
	@mkdir -p $(BENCH_DIR)
	go test $(BENCH_FLAGS) -tags=bench_remy | tee $(BENCH_DIR)/$(REMY_FILE)

bench-do:
	@mkdir -p $(BENCH_DIR)
	go test $(BENCH_FLAGS) -tags=bench_do | tee $(BENCH_DIR)/$(DO_FILE)

bench-dig:
	@mkdir -p $(BENCH_DIR)
	go test $(BENCH_FLAGS) -tags=bench_dig | tee $(BENCH_DIR)/$(DIG_FILE)

benchstat: benchstat-check bench-files
	cd $(BENCH_DIR) && benchstat dig.txt do.txt remy.txt | tee benchstat.txt

update-readme:
	go run ./cmd/update-readme \
		-readme README.md \
		$(BENCH_DIR)/$(DIG_FILE) \
		$(BENCH_DIR)/$(DO_FILE) \
		$(BENCH_DIR)/$(REMY_FILE)

bench: benchstat update-readme

clean:
	rm -rf $(BENCH_DIR)
