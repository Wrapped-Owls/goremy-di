.PHONY: generate test bench

generate:
	go generate ./...

test:
	go test -bench=. -benchmem ./tests

bench: test
