bin_dir=$(shell pwd)/bin


.PHONY: test
test:
	go test -short -v ./...
	golangci-lint run ./...

.PHONY: test-race
test-race:
	go test -short -race ./...

