.PHONY: test test-race lint clean

test:
	go test ./... -v

test-race:
	go test ./... -race -v

lint:
	golangci-lint run ./...

clean:
	go clean -cache -testcache
