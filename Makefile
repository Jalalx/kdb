GIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitHash=$(GIT_HASH)"

.PHONY: clean build run

clean:
	rm -f bin/kdb
	rm -f bin/ask
run:
	go run *.go

build:
	mkdir -p bin/
	go build $(LDFLAGS) -o bin/kdb