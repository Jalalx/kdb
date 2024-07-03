GIT_HASH=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags --abbrev=0 || echo "v0.0.0")
# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitHash=$(GIT_HASH)"


clean:
	rm -f db/knowledgebase.ddb db/knowledgebase.ddb.wal
	rm -f kdb
run:
	go run *.go

build:
	go build $(LDFLAGS) -o bin/kdb