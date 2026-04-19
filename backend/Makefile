.PHONY: build-fw run-lw vet fmt lint test

build-fw:
	go build -o bin/api ./cmd/api

run-fw:
	go run ./cmd/api

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

lint:
	golangci-lint run