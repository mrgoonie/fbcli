BINARY=fbcli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: build clean test lint install

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/fbcli

clean:
	rm -f $(BINARY)

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

install:
	go install $(LDFLAGS) ./cmd/fbcli

fmt:
	go fmt ./...

vet:
	go vet ./...
